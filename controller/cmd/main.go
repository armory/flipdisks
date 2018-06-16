package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/armory/flipdisks/controller/pkg/fontmap"
	"github.com/kevinawoo/flipdots/panel"
	"github.com/nlopes/slack"
)

type MetadataType struct {
	TallerCharacters map[string]int `json:"tallerCharacters"`
	AverageHeight    int            `json:"averageHeight"`
	AverageWidth     int            `json:"averageWidth"`
}

type Row []int
type Letter []Row
type CharmapType map[string]Letter

type Font struct {
	Name     string       `json:"name"`
	Metadata MetadataType `json:"metadata"`
	Charmap  CharmapType  `json:"charmap"`
}

type Board []Row
type Frame map[string]Board

type FlipdiskVideo struct {
	Name         string
	FPS          int        `json:"fps"` // how do we go lower?
	Looping      bool       `json:"looping"`
	Layout       [][]string `json:"layout"`
	Frames       []Board    `json:"frames"`
	SetNullTo    int
	FrameDelayMs int
}

type PanelInfo struct {
	PanelHeight              int
	PanelWidth               int
	PhysicallyDisplayedWidth int
}

type Playlist struct {
	Location             string          `json:"location"`
	Name                 string          `json:"name"`
	Videos               []FlipdiskVideo `json:"videos"`
	Looping              bool            `json:"looping"`
	PanelInfo            PanelInfo
	PanelAddressesLayout [][]int
}

func main() {
	playlist := &Playlist{
		Name:     "demo",
		Location: "armorywall",
		PanelInfo: PanelInfo{
			// actual panels
			PanelWidth:  28,
			PanelHeight: 7,

			PhysicallyDisplayedWidth: 7,
		},
		PanelAddressesLayout: [][]int{
			// actual layouts
			{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			{10, 11, 12, 13, 14, 15, 16, 17, 18, 19},
		},
		Videos: []FlipdiskVideo{
			{
				Name:         "on off",
				Looping:      true,
				FPS:          1,
				SetNullTo:    0,
				FrameDelayMs: 1000,
				Frames:       []Board{{}},
			},
		},
	}

	playlistJson, _ := json.Marshal(playlist)
	playlistJson = playlistJson

	width := flag.Int("w", 28, "width of panel")
	height := flag.Int("h", 7, "width of panel")
	port := flag.String("p", "/dev/tty.SLAB_USBtoUART", "the serial port, empty string to simulate")
	baud := flag.Int("b", 9600, "baud rate of port")

	slackToken := flag.String("slack-token", "", "Go get a slack token")
	flag.Parse()

	width = width
	height = height
	port = port
	port = port
	baud = baud

	var panels [][]*panel.Panel
	for y, row := range playlist.PanelAddressesLayout {
		panels = append(panels, []*panel.Panel{})

		for _, panelAddress := range row {
			p := panel.NewPanel(playlist.PanelInfo.PanelWidth, playlist.PanelInfo.PanelHeight, *port, *baud)
			p.Address = []byte{byte(panelAddress)}

			panels[y] = append(panels[y], p)
			defer p.Close()
		}
	}

	messages := make(chan string)
	go startSlackListener(*slackToken, playlist, panels, messages)
	var msgCharsAsDots []fontmap.Letter

	var board []Row

	for msg := range messages {
		if msg == "debug all panels" || msg == "debug panels" {
			debugPanelAddressByGoingInOrder(panels)
		}

		if strings.Contains(msg, "debug panel") {
			panelAddress, _ := strconv.Atoi(strings.Replace(msg, "debug panel ", "", -1))
			debugSinglePanel(panels, panelAddress)
		}

		// clear the message and the board, ready for the next message
		msgCharsAsDots = msgCharsAsDots[:0]
		board = board[:0]

		// replace slack tokens that are rendered to characters
		msg = strings.Replace(msg, "&lt;", "<", -1)
		msg = strings.Replace(msg, "&gt;", ">", -1)

		msgCharsAsDots = fontmap.Render(msg)

		// we have to convert our long array of dotCharacters to a virtual board
		var longestLine, lineNumber int
		longestLine = 0
		lineNumber = 0
		lineMaxWidth := playlist.PanelInfo.PhysicallyDisplayedWidth * len(playlist.PanelAddressesLayout[0])

		// join the letters together to form one long string
		for charIndexInMessage, charAsDots := range msgCharsAsDots {
			if longestLine+len(charAsDots[0]) > lineMaxWidth || msg[charIndexInMessage] == '\n' {
				lineNumber++
				longestLine = 0

				if msg[charIndexInMessage] == '\n' && charAsDots == nil {
					continue
				}
			}

			for charRowIndex, charRow := range charAsDots {
				boardCharRowIndex := charRowIndex + (lineNumber * fontmap.TI84.Metadata.MaxHeight)
				//log.Println(len(board), boardCharRowIndex, charRowIndex, lineNumber, fontmap.TI84.Metadata.MaxHeight)
				for len(board) <= boardCharRowIndex { // 2 < 2
					board = append(board, Row{})
				}

				//log.Printf("writing to line number: %d, boardCharRowIndex %d", lineNumber, boardCharRowIndex)
				//log.Print("row:", charRow)
				board[boardCharRowIndex] = append(board[boardCharRowIndex], charRow...)

				if longestLine < len(board[boardCharRowIndex]) {
					longestLine = len(board[boardCharRowIndex])
				}
			}
		}

		printBoard(board)

		frameIndex := 0
		frameIndex = frameIndex

		// convert virtual board to a physical board
		for x := 0; x < len(board); x++ {
			for y := 0; y < len(board[x]); y++ {
				panelXCoord := x / playlist.PanelInfo.PanelWidth
				panelYCoord := y / playlist.PanelInfo.PanelHeight

				if panelXCoord >= len(playlist.PanelAddressesLayout) {
					log.Printf("Warning: Frame %d row %d, exceeds specified HEIGHT %d, dropping the rest of it.", frameIndex, y, playlist.PanelInfo.PanelWidth)
					continue
				}

				if panelYCoord >= len(playlist.PanelAddressesLayout[panelXCoord]) {
					log.Printf("Warning: Frame %d cell(%d,%d) exceeds specified WIDTH %d, dropping the rest of it.", frameIndex, x, y, playlist.PanelInfo.PanelWidth)
					continue
				}

				p := panels[panelXCoord][panelYCoord]

				// which dot should we set?
				dotXCoord := x % playlist.PanelInfo.PanelWidth
				dotYCoord := y % playlist.PanelInfo.PanelHeight
				dotValue := board[x][y] == 1
				//log.Printf("Setting panel(%d,%d), adddress %d, dot(%d,%d) with %t", panelXCoord, panelYCoord, p.Address, dotXCoord, dotYCoord, dotValue)
				p.Set(dotXCoord, dotYCoord, dotValue)
			}
		}

		// print the panels
		for y, row := range panels {
			for x, p := range row {
				y = y
				x = x
				//fmt.Println(x, y, p.Address)
				//p.PrintState()
				p.Send()
				p.Clear(false)
				//p.Send()
			}
		}
	}
}

func debugPanelAddressByGoingInOrder(panels [][]*panel.Panel) {
	// clear all boards
	for _, row := range panels {
		for _, p := range row {
			p.Clear(false)
			p.Send()
		}
	}

	dotState := false
	for {
		dotState = !dotState

		for y, row := range panels {
			for x, p := range row {
				x = x
				y = y
				//if p.Address[0] == byte(1) {
				//fmt.Println(x, y, p.Address, dotState)
				p.Clear(dotState)
				p.Send()
				time.Sleep(time.Duration(250) * time.Millisecond)
				//}
				//p.Send()
			}
		}
	}
}

func debugSinglePanel(panels [][]*panel.Panel, address int) {
	// clear all boards
	for _, row := range panels {
		for _, p := range row {
			p.Clear(false)
			p.Send()
		}
	}

	dotState := false
	for {
		dotState = !dotState

		for y, row := range panels {
			for x, p := range row {
				x = x
				y = y
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

func printBoard(board Board) {
	for x := 0; x < len(board); x++ {
		line := ""
		for y := 0; y < len(board[x]); y++ {
			if board[x][y] == 1 {
				line += "⚫️"
			} else {
				line += "⚪️"
			}
		}
		log.Println(line)
	}
}


func startVideoPlayer(playlist *Playlist, panels [][]*panel.Panel) {
	for _, video := range playlist.Videos {
		for {
			for frameIndex, frame := range video.Frames {
				for y, row := range frame {
					panelRow := y / playlist.PanelInfo.PanelHeight

					if panelRow >= len(playlist.PanelAddressesLayout) {
						log.Printf("Warning: Frame %d row %d, exceeds specified HEIGHT %d, dropping the rest of it.", frameIndex, y, playlist.PanelInfo.PanelWidth)
						break
					}
					for x, cellValue := range row {
						panelColumn := x / playlist.PanelInfo.PanelWidth

						if panelColumn >= len(playlist.PanelAddressesLayout[panelRow]) {
							log.Printf("Warning: Frame %d cell(%d,%d) exceeds specified WIDTH %d, dropping the rest of it.", frameIndex, x, y, playlist.PanelInfo.PanelWidth)
							break
						}

						p := panels[panelRow][panelColumn]
						p.Set(x%playlist.PanelInfo.PanelWidth, y%playlist.PanelInfo.PanelHeight, cellValue == 1)
					}
				}

				for y, row := range panels {
					for x, p := range row {
						fmt.Println(x, y, p.Address)
						p.PrintState()
						p.Send()
					}
				}
				time.Sleep(time.Duration(video.FrameDelayMs) * time.Millisecond)
			}

			if !video.Looping {
				break
			}
		}
	}
}

func startSlackListener(slackToken string, playlist *Playlist, panels [][]*panel.Panel, messages chan string) {
	api := slack.New(slackToken)
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	api.SetDebug(false)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		//fmt.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			//fmt.Printf("Message: %+v\n\n", ev)
			fmt.Printf("GOT MESSAGE: %+v\n", ev.Msg.Text)
			flipbotUserId := "UASEXQA04"

			if ev.SubMessage != nil {
				// someone edited their old message, let's display it
				messages <- ev.SubMessage.Text
			} else {
				messages <- ev.Msg.Text
			}

			if strings.Contains(ev.Msg.Text, flipbotUserId) {
				flipTableWordList := []string{
					"flip",
					"table",
				}

				curseWords := []string{
					"fuck",
					"god",
					"damn",
					"ass",
				}

				for _, word := range curseWords {
					if strings.Contains(ev.Msg.Text, word) {
						rtm.SendMessage(rtm.NewOutgoingMessage("Yo, watch your language, you dick head...", "DAZ6XPGJ1"))
						break
					}
				}

				for _, word := range flipTableWordList {
					if strings.Contains(ev.Msg.Text, word) {
						rtm.SendMessage(rtm.NewOutgoingMessage("Flipping the table (╯°□°）╯︵ ┻━┻", "DAZ6XPGJ1"))
						break
					}
				}
			}

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return

		default:
		}
	}
}
