package ngrok

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Config struct {
	Host          string
	openPorts     []Port
	tunnelTimeout time.Duration
	activeWG      *sync.WaitGroup
}

func New() *Config {
	return &Config{
		Host:          "http://localhost:4040",
		tunnelTimeout: time.Minute * 15,
		activeWG:      &sync.WaitGroup{},
	}
}

type Protocol string

const (
	ProtocolTCP  = Protocol("tcp")
	ProtocolHTTP = Protocol("http")
)

type Port struct {
	Addr      string `json:"addr"`
	Proto     string `json:"proto"`
	Name      string `json:"name"`
	PublicURL string
}

type ngrokPortRes struct {
	Name      string `json:"name"`
	URI       string `json:"uri"`
	PublicURL string `json:"public_url"`
	Proto     string `json:"proto"`
	Config    struct {
		Addr    string `json:"addr"`
		Inspect bool   `json:"inspect"`
	} `json:"config"`
}

func (c *Config) GetTimeout() time.Duration {
	return c.tunnelTimeout
}
func (c *Config) IsDaemonActive() bool {
	_, err := c.ListTunnels()
	if err != nil {
		return false
	}
	return true
}

func (c *Config) StartTunnel(name string, proto Protocol, port int) (Port, error) {
	p := struct {
		Name  string `json:"name"`
		Addr  string `json:"addr"`
		Proto string `json:"proto"`
	}{
		Name:  name,
		Addr:  strconv.Itoa(port),
		Proto: string(proto),
	}

	fmt.Println(fmt.Sprintf("starting tunnel '%s' %s %s'", p.Name, p.Proto, p.Addr))

	s, _ := json.Marshal(p)
	res, err := http.Post(c.Host+"/api/tunnels", "application/json", bytes.NewBuffer(s))
	if err != nil {
		return Port{}, err
	}
	defer res.Body.Close()

	s, _ = ioutil.ReadAll(res.Body)

	var portRes ngrokPortRes
	_ = json.Unmarshal(s, &portRes)

	createdPort := Port{
		Addr:      portRes.Config.Addr,
		Proto:     portRes.Proto,
		Name:      portRes.Name,
		PublicURL: portRes.PublicURL,
	}

	c.openPorts = append(c.openPorts, createdPort)
	c.activeWG.Add(1)

	fmt.Println(fmt.Sprintf("succefully started tunnel '%s' %s %s', tunnel will expire in %s", p.Name, p.Proto, p.Addr, c.tunnelTimeout.String()))
	go func() {
		t := time.After(c.tunnelTimeout)
		<-t
		_ = c.StopAllTunnels()
	}()

	return createdPort, nil
}

// StopAllTunnels intends to stop everything, even ones that we don't manage.
// This is to make sure that we're not exposing anything we don't want to expose
func (c *Config) StopAllTunnels() []error {
	tunnels, err := c.ListTunnels()
	if err != nil {
		return []error{err}
	}

	var errs []error
	for _, t := range tunnels {
		fmt.Println(fmt.Sprintf("timeout: closing '%s' %s %s'", t.Name, t.Proto, t.Addr))
		req, err := http.NewRequest(http.MethodDelete, c.Host+"/api/tunnels/"+t.Name, nil)
		if err != nil {
			fmt.Println(fmt.Sprintf("could not stop tunnel '%s' %s %s'", t.Name, t.Proto, t.Addr))
			errs = append(errs, err)
		}

		client := &http.Client{}
		res, _ := client.Do(req)
		defer res.Body.Close()

		for _, p := range c.openPorts {
			if p.Name == t.Name {
				c.activeWG.Done()
			}
		}

		fmt.Println(fmt.Sprintf("succesfully stopped tunnel '%s' %s %s'", t.Name, t.Proto, t.Addr))
	}

	return errs
}

func (c *Config) Active() *sync.WaitGroup {
	return c.activeWG
}

func (c *Config) ListTunnels() ([]Port, error) {
	res, err := http.Get(c.Host + "/api/tunnels")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var b struct {
		Tunnels []struct {
			Name   string
			Proto  string
			Config struct {
				Addr string
			}
		}
	}

	s, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(s, &b)
	if err != nil {
		return nil, err
	}

	var tunnels []Port
	for _, t := range b.Tunnels {
		tunnels = append(tunnels, Port{
			Addr:  strings.Split(t.Config.Addr, ":")[1],
			Proto: t.Proto,
			Name:  t.Name,
		})
	}

	return tunnels, nil
}
