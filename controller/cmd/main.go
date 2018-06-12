package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
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
	PanelHeight int
	PanelWidth  int
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
			// debug panels
			//PanelWidth: 28,
			//PanelHeight: 100,

			// actual panels
			PanelWidth:  7,
			PanelHeight: 28,
		},
		PanelAddressesLayout: [][]int{
			// debug layouts
			//{1, 2, 3, 4, 5,},
			//{6, 7, 8, 9, 10},
			//{1},

			// actual layouts
			{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			{11, 12, 13, 14, 15, 16, 17, 18, 19, 20},

		},
		Videos: []FlipdiskVideo{
			{
				Name:         "on off",
				Looping:      true,
				FPS:          1,
				SetNullTo:    0,
				FrameDelayMs: 1000,
				Frames:       []Board{{},},
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
	//go func() { messages <- "Good morning" }()
	//lastSeenMessage := ""
	var msgCharsAsDots []fontmap.Letter

	var board []Row

	for msg := range messages {
		msgCharsAsDots = msgCharsAsDots[:0]
		board = board[:0]

		/*
		if the message is "hello world"
		then then we'll create ["h","e","l","l","o"," ","w","o","r","l","d"]
		which then will have each character turned into a 2x2 matrix of dots
		the final output will by an array of 2x2 matrixes
		 */
		msg = strings.Replace(msg, "&lt;", "<", -1)
		msg = strings.Replace(msg, "&gt;", ">", -1)


		for _, char := range strings.Split(msg, "") {
			if char == " " {
				msgCharsAsDots = append(msgCharsAsDots, generateSpace(2, fontmap.TI84.Metadata.MaxHeight, 0)) // random magic 2 for pretty printing letters with tails
			} else {
				if dotLetter, charExists := fontmap.TI84.Charmap[char]; charExists {
					msgCharsAsDots = append(msgCharsAsDots, addKerning(dotLetter, 0))
				} else {
					msgCharsAsDots = append(msgCharsAsDots, generateSpace(4, fontmap.TI84.Metadata.MaxHeight, 1))
				}
			}
		}

		var longestLine, lineNumber int
		longestLine = 0
		lineNumber = 0
		lineMaxWidth := playlist.PanelInfo.PanelWidth * len(playlist.PanelAddressesLayout[0])

		// join the letters together to form one long string
		for _, charAsDots := range msgCharsAsDots {
			for charRowIndex, charRow := range charAsDots {
				boardCharRowIndex := charRowIndex + lineNumber*fontmap.TI84.Metadata.MaxHeight
				if len(board) <= boardCharRowIndex {
					board = append(board, Row{})
				}

				//fmt.Println("writing to line number:", lineNumber)
				board[boardCharRowIndex] = append(board[boardCharRowIndex], charRow...)

				if longestLine < len(board[boardCharRowIndex]) {
					longestLine = len(board[boardCharRowIndex])
				}
			}

			if longestLine >= lineMaxWidth {
				lineNumber++
				board = append(board, Row{})

				// add a whitespace line

				//for i := 0; i < longestLine; i++ {
				//	board[len(board)-1] = append(board[len(board)-1], 0)
				//}
				longestLine = 0
			}
		}

		frameIndex := 0

		for y, row := range board {
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
				p.Clear(false)
			}
		}
	}
}

func generateSpace(width int, height int, value int) fontmap.Letter {
	var space fontmap.Letter
	for j := 0; j < height; j++ {
		var row fontmap.Row
		for i := 0; i < width; i++ {
			row = append(row, value)
		}

		space = append(space, row)
	}

	return space
}

func addKerning(letter fontmap.Letter, amountOfKerning int) fontmap.Letter {
	var kernedLetter fontmap.Letter
	for _, row := range letter {
		kernedRow := row
		for j := 0; j < amountOfKerning; j++ {
			kernedRow = append(kernedRow, 0)
		}
		kernedLetter = append(kernedLetter, kernedRow)
	}

	return kernedLetter
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
			messages <- ev.Msg.Text
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
