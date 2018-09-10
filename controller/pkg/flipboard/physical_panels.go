package flipboard

import (
	"fmt"
	"time"

	"github.com/kevinawoo/flipdots/panel"
	"github.com/sirupsen/logrus"
)

type PanelInfo struct {
	PanelHeight              int
	PanelWidth               int
	PhysicallyDisplayedWidth int
	Port                     string
	Baud                     int
}

type PanelAddress int
type PanelLayout [][]PanelAddress

func CreatePanels(panelInfo PanelInfo, panelLayout PanelLayout) (*[][]panel.Panel, error) {
	var panels [][]panel.Panel

	for y, row := range panelLayout {
		panels = append(panels, []panel.Panel{})

		for _, panelAddress := range row {
			p, err := panel.NewPanel(panelInfo.PanelWidth, panelInfo.PanelHeight, panelInfo.Port, panelInfo.Baud)
			if err != nil {
				return nil, err
			}

			p.Address = []byte{byte(panelAddress)}

			panels[y] = append(panels[y], *p)
		}
	}
	return &panels, nil
}

func (b *Flipboard) DebugPanelAddressByGoingInOrder() {
	// clear all boards
	for _, row := range *b.panels {
		for _, p := range row {
			p.Clear(false)
			p.Send()
		}
	}

	dotState := false
	for {
		select {
		case <-b.newMessage:
			// got a new message, stop debugging
			go func() { b.newMessage <- true }() // let the play function know to continue again
			return
		default:
			dotState = !dotState

			for _, row := range *b.panels {
				for _, p := range row {
					p.Clear(dotState)
					p.Send()
					time.Sleep(time.Duration(250) * time.Millisecond)
				}
			}
		}
	}
}

func (b *Flipboard) DebugSinglePanel(address int) {
	// clear all boards
	for _, row := range *b.panels {
		for _, p := range row {
			p.Clear(false)
			p.Send()
		}
	}

	dotState := false
	for {
		select {
		case <-b.newMessage:
			// got a new message, stop debugging
			go func() { b.newMessage <- true }() // let the play function know to continue again
			return
		default:
			dotState = !dotState

			for y, row := range *b.panels {
				for x, p := range row {
					if p.Address[0] == byte(address) {
						fmt.Println(x, y, p.Address, dotState)
						p.Clear(dotState)
						p.Send()
						time.Sleep(time.Duration(500) * time.Millisecond)
					}
				}
			}
		}
	}
}

func (b *Flipboard) SetAll(val bool) {
	for _, row := range *b.panels {
		for _, p := range row {
			p.Clear(val)
		}
	}
}

func (b *Flipboard) GetPanel(x, y int) (*panel.Panel) {
	panels := *b.panels
	return &panels[x][y]
}

func (b *Flipboard) Send() () {
	for y, row := range *b.panels {
		for x, p := range row {
			//p.PrintState()
			err := p.Send()
			if err != nil {
				logrus.Errorf("could not send to panel (%d,%d): %s", y, x, err)
			}
		}
	}
}
