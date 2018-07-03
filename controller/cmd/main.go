package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/armory/flipdisks/controller/pkg/fontmap"
	"github.com/kevinawoo/flipdots/panel"
	"github.com/nlopes/slack"
	"gopkg.in/yaml.v2"
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

type FlipBoardDisplayOptions struct {
	Append   bool   `yaml:"append"`
	Align    string `yaml:"align"`
	FontSize int    `yaml:"font-size"`
	Kerning  int    `yaml:"kerning"`
}

var flipBoardDisplayOptions FlipBoardDisplayOptions

func main() {
	log.Print("Starting")

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
	flipBoardDisplayOptions.Append = false
	go startSlackListener(*slackToken, playlist, panels, messages)
	var msgCharsAsDots []fontmap.Letter

	var virtualBoard VirtualBoard

	for msg := range messages {
		if msg == "debug all panels" || msg == "debug panels" {
			debugPanelAddressByGoingInOrder(panels)
		}

		if strings.Contains(msg, "debug panel") {
			panelAddress, _ := strconv.Atoi(strings.Replace(msg, "debug panel ", "", -1))
			debugSinglePanel(panels, panelAddress)
		}

		// clear the message and the virtualBoard, ready for the next message
		msgCharsAsDots = msgCharsAsDots[:0]
		virtualBoard = virtualBoard[:0]

		msgCharsAsDots = fontmap.Render(msg)
		virtualBoard = createVirtualBoard(playlist.PanelInfo.PhysicallyDisplayedWidth, len(playlist.PanelAddressesLayout[0]), msgCharsAsDots, msg)

		printBoard(virtualBoard)

		frameIndex := 0
		frameIndex = frameIndex

		// convert virtual virtualBoard to a physical virtualBoard
		for x := 0; x < len(virtualBoard); x++ {
			for y := 0; y < len(virtualBoard[x]); y++ {
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
				dotValue := virtualBoard[x][y] == 1
				//log.Printf("Setting panel(%d,%d), adddress %d, dot(%d,%d) with %t", panelXCoord, panelYCoord, p.Address, dotXCoord, dotYCoord, dotValue)
				p.Set(dotXCoord, dotYCoord, dotValue)
			}
		}

		// send our virtual panels to the physical virtualBoard
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

type VirtualBoard []fontmap.Row

func createVirtualBoard(panelWidth int, numberOfPanelsWide int, msgCharsAsDots []fontmap.Letter, msg string) VirtualBoard {
	// we have to convert our long array of dotCharacters to a virtual board
	var longestLine, lineNumber int
	longestLine = 0
	lineNumber = 0
	lineMaxWidth := panelWidth * numberOfPanelsWide
	var virtualBoard VirtualBoard

	// join the letters together to form one long string
	for charIndexInMessage := 0; charIndexInMessage < len(msgCharsAsDots); charIndexInMessage++ {
		charAsDots := msgCharsAsDots[charIndexInMessage]

		// handle line breaks
		if msg[charIndexInMessage] == '\n' && charAsDots == nil {
			lineNumber++
			longestLine = 0
			continue
		}

		// try to word break, if it doesn't work, then we'lll need to character break
		if msg[charIndexInMessage] == ' ' {
			unprocessedStringMsg := msg[charIndexInMessage:] // msg will look something like: "   bbb"
			unprocessedDotMessage := msgCharsAsDots[charIndexInMessage:]

			matchPos := regexp.MustCompile(`\S+`).FindStringIndex(unprocessedStringMsg) // matchPos[0] will be the first "b"
			nextDotWord := unprocessedDotMessage[matchPos[0]:matchPos[1]]

			// find the width of dots for the word
			wordDotWidth := 0
			for _, dotWord := range nextDotWord {
				wordDotWidth += len(dotWord[0])
			}

			// since we're breaking on the word, we should discard all the whitespace before the word
			if longestLine+wordDotWidth > lineMaxWidth {
				lineNumber++
				longestLine = 0

				// advance our pointer to the beginning of the next word
				charIndexInMessage += matchPos[0]
				charAsDots = msgCharsAsDots[charIndexInMessage]
			}
		} else if longestLine+len(charAsDots[0]) > lineMaxWidth {
			// if there's no spaces, and the word is super long, let's fallback and do a character break
			lineNumber++
			longestLine = 0
		}

		// write character to the virtual board
		for charRowIndex, charRow := range charAsDots {
			boardCharRowIndex := charRowIndex + (lineNumber * fontmap.TI84.Metadata.MaxHeight)

			// create all missing rows from the virtual board, up to our current boardCharRowIndex
			for len(virtualBoard) <= boardCharRowIndex {
				virtualBoard = append(virtualBoard, fontmap.Row{})
			}

			virtualBoard[boardCharRowIndex] = append(virtualBoard[boardCharRowIndex], charRow...)

			// keep track of the longest char row for the line
			if longestLine < len(virtualBoard[boardCharRowIndex]) {
				longestLine = len(virtualBoard[boardCharRowIndex])
			}

		}
	}
	return virtualBoard
}

func (board VirtualBoard) String() string {
	line := ""
	for x := 0; x < len(board); x++ {
		for y := 0; y < len(board[x]); y++ {
			if board[x][y] == 1 {
				line += "⚫️"
			} else {
				line += "⚪️"
			}
		}
		line += "\n"
	}

	return line
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

func printBoard(board VirtualBoard) {
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

func startSlackListener(slackToken string, playlist *Playlist, panels [][]*panel.Panel, flipboardMsgChn chan string) {
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
			handleSlackMsg(ev, rtm, flipboardMsgChn)

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return

		default:
		}
	}
}

func handleSlackMsg(ev *slack.MessageEvent, rtm *slack.RTM, flipboardMsgChn chan string) {
	rawMsg := ev.Msg.Text
	fmt.Printf("Raw Slack Message: %+v\n", rawMsg)

	msgOptions := regexp.MustCompile("\\s*---\\s*").Split(rawMsg, -1)
	msg := msgOptions[0]
	if len(msgOptions) > 1 {
		setFlipboardOptions(msgOptions[1])
	}

	msg = renderSlackUsernames(msg, rtm)
	msg = cleanupSlackEncodedCharacters(msg)

	fmt.Printf("Rendering message: %+v\n", msg)
	if msg == "help" {
		respondWithHelpMsg(rtm, ev.Msg.Channel)
		return
	}

	if ev.SubMessage != nil {
		// someone edited their old message, let's display it
		flipboardMsgChn <- ev.SubMessage.Text
	} else {
		flipboardMsgChn <- msg
	}
}

func cleanupSlackEncodedCharacters(msg string) string {
	// replace slack tokens that are rendered to characters
	msg = strings.Replace(msg, "&lt;", "<", -1)
	msg = strings.Replace(msg, "&gt;", ">", -1)
	return msg
}

func setFlipboardOptions(rawOptions string) {
	err := yaml.Unmarshal([]byte(rawOptions), &flipBoardDisplayOptions)
	err = err
	fmt.Printf("%#v \n", flipBoardDisplayOptions)
}

func renderSlackUsernames(msg string, rtm *slack.RTM) string {
	userIds := regexp.MustCompile("<@\\w+>").FindAllString(msg, -1)
	for _, slackFmtMsgUserId := range userIds {
		// in the message we'll receive something like "<@U123123>", the id is actually "U123123"
		userId := strings.Replace(strings.Replace(slackFmtMsgUserId, "<@", "", 1), ">", "", 1)

		user, err := rtm.GetUserInfo(userId)
		if err != nil {
			name := user.Name
			if user.Profile.FirstName != "" {
				name = user.Profile.FirstName
			}
			msg = strings.Replace(msg, "<@"+user.ID+">", name, -1)
		}
	}
	return msg
}

func respondWithHelpMsg(rtm *slack.RTM, channelId string) {
	msg := `Send me a DM and I'll display that. 
You can also change settings by doing:
`

	msg += "```"
	msg += `
Your cool message here!
--- 
append: true/false	 // overwrite or add to the board
`

// we would like to add support for this in the future
//align: center center   // horizontal vertical
//kerning: 0	         // spacing between letters
//font-size: 1           // ??


	msg += "```"
	rtm.SendMessage(rtm.NewOutgoingMessage(msg, channelId))
}
