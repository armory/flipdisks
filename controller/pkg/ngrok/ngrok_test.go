package ngrok

import (
	"testing"
	"time"
)

func TestIntegration(t *testing.T) {
	c := New()
	c.tunnelTimeout = time.Second * 10
	c.Host= "http://localhost:4040"

	c.StartTunnel("randoPort", ProtocolTCP, 4958)

	c.Active().Wait()
}
