package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/armory/flipdisks/controller/pkg/image"
	"github.com/armory/flipdisks/controller/pkg/fontmap"
	"github.com/armory/flipdisks/controller/pkg/github"
	"github.com/armory/flipdisks/controller/pkg/options"
	"github.com/armory/flipdisks/controller/pkg/slackbot"
	"github.com/kevinawoo/flipdots/panel"
	log "github.com/sirupsen/logrus"
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

var githubToken *string
var countdownDate string

func main() {
	log.Print("Starting flipdisk controller")

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

	//width := flag.Int("w", 28, "width of panel")
	//height := flag.Int("h", 7, "width of panel")
	port := flag.String("p", "/dev/tty.SLAB_USBtoUART", "the serial port, empty string to simulate")
	baud := flag.Int("b", 9600, "baud rate of port")

	slackToken := flag.String("slack-token", "", "Go get a slack token")
	githubToken = flag.String("github-token", "", "Go get a github token")
	flag.StringVar(&countdownDate, "countdown", "", fmt.Sprintf("Specify the countdown date in YYYY-MM-DD format"))
	flag.Parse()

	g, err := github.New(github.Token(*githubToken))
	if err != nil {
		log.Error("Could not create githubClient, hopefully everything will work!")
	}
	githubEmojiLookup, err := g.GetEmojis()
	if err != nil {
		log.Error("Could not get emojis from Github", err)
	}

	panels, err := createPanels(playlist, port, baud)
	if err != nil {
		log.Fatal(err)
	}

	msgsChan := make(chan options.FlipBoardDisplayOptions)

	_ = slackToken
	_ = msgsChan
	slack := slackbot.NewSlack(countdownDate, githubEmojiLookup)
	go slack.StartSlackListener(*slackToken, msgsChan)

	for msg := range msgsChan {
		DisplayMessageToPanels(msg, panels, playlist)
	}
}

func createPanels(playlist *Playlist, port *string, baud *int) ([][]*panel.Panel, error) {
	var panels [][]*panel.Panel

	for y, row := range playlist.PanelAddressesLayout {
		panels = append(panels, []*panel.Panel{})

		for _, panelAddress := range row {
			p, err := panel.NewPanel(playlist.PanelInfo.PanelWidth, playlist.PanelInfo.PanelHeight, *port, *baud)
			if err != nil {
				return nil, err
			}

			p.Address = []byte{byte(panelAddress)}

			panels[y] = append(panels[y], p)
		}
	}
	return panels, nil
}

func DisplayMessageToPanels(msg options.FlipBoardDisplayOptions, panels [][]*panel.Panel, playlist *Playlist) {
	if msg.Message == "debug all panels" || msg.Message == "debug panels" {
		debugPanelAddressByGoingInOrder(panels)
	}
	if strings.Contains(msg.Message, "debug panel") {
		panelAddress, _ := strconv.Atoi(strings.Replace(msg.Message, "debug panel ", "", -1))
		debugSinglePanel(panels, panelAddress)
	}

	virtualBoard := renderVirtualBoard(msg, playlist)

	frameIndex := 0
	frameIndex = frameIndex
	fill := msg.Fill == "true"

	// if no fill is provided, let's try to set autofill
	if msg.Fill == "" {
		var sum int

		// Go across the top to add up all the values
		for x := range virtualBoard[0] {
			sum += virtualBoard[0][x]
		}

		// go across the bottom to add up all the values
		for x := range virtualBoard[len(virtualBoard)-1] {
			sum += virtualBoard[len(virtualBoard)-1][x]
		}

		// go on the left and right side to add up all the values
		for _, row := range virtualBoard {
			// sometimes the row will be empty, because of a \n, let's just ignore it
			if len(row) > 0 {
				sum += row[0] // left y going down
			}

			// if for some reason it's just a single row, we'll have already taken care of adding the sum before
			if len(row) > 1 {
				sum += row[len(row)-1] // right y going down
			}
		}

		height := len(virtualBoard)
		width := len(virtualBoard[0])
		fill = float32(sum)/float32(2*(width+height)) >= .5 // magic number
		//fmt.Println("setting autofill to be: ", fill)
	}
	// set the fill value
	for _, row := range panels {
		for _, p := range row {
			p.Clear(fill)
		}
	}

	// set alignment options
	msg.XAlign, msg.YAlign = options.GetAlignOptions(msg.Align)

	printBoard(virtualBoard)

	// the library flipped height and width by accident, we'll work around it
	panelWidth := playlist.PanelInfo.PanelHeight
	panelHeight := playlist.PanelInfo.PanelWidth

	// convert virtual virtualBoard to a physical virtualBoard
	boardWidth := panelWidth * len(playlist.PanelAddressesLayout[0])
	boardHeight := panelHeight * len(playlist.PanelAddressesLayout)
	xOffSet, yOffSet := findOffSets(msg, virtualBoard, boardWidth, boardHeight)
	for y := 0; y < len(virtualBoard); y++ {
		for x := 0; x < len(virtualBoard[y]); x++ {
			// which dot should we set?
			panelXCoord := (x + xOffSet) / panelWidth
			panelYCoord := (y + yOffSet) / panelHeight
			dotXCoord := (x + xOffSet) % panelWidth
			dotYCoord := (y + yOffSet) % panelHeight

			if dotXCoord < 0 || dotYCoord < 0 || panelXCoord < 0 || panelYCoord < 0 {
				continue
			}

			if panelYCoord >= len(playlist.PanelAddressesLayout) {
				log.Printf("Warning: Frame %d row %d, exceeds specified HEIGHT %d, dropping the rest of it.", frameIndex, x, panelHeight)
				continue
			}

			if panelXCoord >= len(playlist.PanelAddressesLayout[panelYCoord]) {
				log.Printf("Warning: Frame %d cell(%d,%d) exceeds specified WIDTH %d, dropping the rest of it.", frameIndex, y, x, panelWidth)
				continue
			}

			//log.Printf("Setting panel(%d,%d), adddress %d, dot(%d,%d) with %t", panelYCoord, panelXCoord, p.Address, dotYCoord, dotXCoord, dotValue)

			// there's a bug in this library, where x and y are flipped. we need to handle this later
			p := panels[panelYCoord][panelXCoord]
			dotValue := virtualBoard[y][x] == 1
			p.Set(dotYCoord, dotXCoord, dotValue)
		}
	}
	// send our virtual panels to the physical virtualBoard
	for y, row := range panels {
		for x, p := range row {
			//p.PrintState()
			err := p.Send()
			if err != nil {
				log.Errorf("could not send to panel (%d,%d): %s", y, x, err)
			}
		}
	}
}

type VirtualBoardCache map[options.FlipBoardDisplayOptions]VirtualBoard

var virtualBoardCache VirtualBoardCache

func renderVirtualBoard(msg options.FlipBoardDisplayOptions, playlist *Playlist) VirtualBoard {
	var virtualBoard VirtualBoard

	// try returning the cache
	if virtualBoardCache == nil {
		virtualBoardCache = VirtualBoardCache{}
	} else if virtualBoardCache[msg] != nil {
		return virtualBoardCache[msg]
	}

	matchedUrls := regexp.MustCompile("http.?://.*.(png|jpe?g)").FindStringSubmatch(msg.Message)
	if len(matchedUrls) > 0 {
		maxWidth := uint(playlist.PanelInfo.PanelHeight * len(playlist.PanelAddressesLayout[0]))
		maxHeight := uint(playlist.PanelInfo.PanelWidth * len(playlist.PanelAddressesLayout))
		virtualBoard = image.Download(maxWidth, maxHeight, matchedUrls[0], msg.Inverted, msg.BWThreshold)
	} else {
		msgCharsAsDots := fontmap.Render(msg.Message)
		virtualBoard = createVirtualBoard(playlist.PanelInfo.PhysicallyDisplayedWidth, len(playlist.PanelAddressesLayout[0]), msgCharsAsDots, msg.Message)

		// todo, it would be nice to just invert it without through the whole board again
		// handle inverting for words
		if msg.Inverted {
			for _, row := range virtualBoard {
				for charIndex, x := range row {
					if x == 0 {
						row[charIndex] = 1
					} else {
						row[charIndex] = 0
					}
				}
			}
		}
	}

	// let's cache the result
	virtualBoardCache[msg] = virtualBoard
	return virtualBoard
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
			for _, dotChar := range nextDotWord {
				if len(dotChar) > 0 {
					wordDotWidth += len(dotChar[0])
				}
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

func findOffSets(options options.FlipBoardDisplayOptions, virtualBoard VirtualBoard, boardWidth, boardHeight int) (int, int) {
	var xOffSet int

	switch options.XAlign {
	case "left":
		xOffSet = 0
	case "center":
		xOffSet = (boardWidth - len(virtualBoard[0])) / 2
		fmt.Println(xOffSet)
	case "right":
		xOffSet = boardWidth - len(virtualBoard[0])
	default:
		xOffSet, _ = strconv.Atoi(options.XAlign)
	}

	var yOffSet int

	switch options.YAlign {
	case "top":
		// we don't do anything
	case "center":
		yOffSet = (boardHeight - len(virtualBoard)) / 2
	case "bottom":
		yOffSet = boardHeight - len(virtualBoard)
	default:
		yOffSet, _ = strconv.Atoi(options.YAlign)
	}

	return xOffSet, yOffSet
}
