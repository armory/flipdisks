package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/armory/flipdisks/controller/pkg/fontmap"
	"github.com/armory/flipdisks/controller/pkg/github"
	"github.com/kevinawoo/flipdots/panel"
	"github.com/nfnt/resize"
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
	Message     string `yaml:"message"`
	DisplayTime int    `yaml:"displayTime"`
	Append      bool   `yaml:"append"`
	Align       string `yaml:"align"`
	xAlign      string
	yAlign      string
	FontSize    int    `yaml:"font-size"`
	Kerning     int    `yaml:"kerning"`
	Inverted    bool   `yaml:"inverted"`
	BWThreshold int    `yaml:"bwThreshold"`
	Fill        string `yml:"fill"`
}

var githubToken *string

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
	flag.Parse()

	g, err := github.New(github.Token(*githubToken))
	if err != nil {
		log.Panic("Could not create githubClient")
	}
	githubEmojiLookup, err = g.GetEmojis()
	if err != nil {
		log.Panicln("Could not get emojis from Github", err)
	}

	panels := createPanels(playlist, port, baud)

	msgsChan := make(chan FlipBoardDisplayOptions)

	_ = slackToken
	_ = msgsChan
	go startSlackListener(*slackToken, msgsChan)

	for msg := range msgsChan {
		DoWOrkFigureOutAName(msg, panels, playlist)
	}
}

func createPanels(playlist *Playlist, port *string, baud *int) [][]*panel.Panel {
	var panels [][]*panel.Panel

	for y, row := range playlist.PanelAddressesLayout {
		panels = append(panels, []*panel.Panel{})

		for _, panelAddress := range row {
			p := panel.NewPanel(playlist.PanelInfo.PanelWidth, playlist.PanelInfo.PanelHeight, *port, *baud)
			p.Address = []byte{byte(panelAddress)}

			panels[y] = append(panels[y], p)
		}
	}
	return panels
}

func DoWOrkFigureOutAName(msg FlipBoardDisplayOptions, panels [][]*panel.Panel, playlist *Playlist) {
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
	// if autofill, try to determine the average around the borders and use that
	fill := msg.Fill == "true"
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
	printBoard(virtualBoard)
	// the library fliped height and width by accident, we'll work around it
	panelWidth := playlist.PanelInfo.PanelHeight
	panelHeight := playlist.PanelInfo.PanelWidth
	// god damn it, its really confusing
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
	for _, row := range panels {
		for _, p := range row {
			//p.PrintState()
			p.Send()
		}
	}
}

type VirtualBoardCache map[FlipBoardDisplayOptions]VirtualBoard

var virtualBoardCache VirtualBoardCache

func renderVirtualBoard(msg FlipBoardDisplayOptions, playlist *Playlist) VirtualBoard {
	var virtualBoard VirtualBoard

	// try returning the cache
	if virtualBoardCache == nil {
		virtualBoardCache = VirtualBoardCache{}
	} else if virtualBoardCache[msg] != nil {
		return virtualBoardCache[msg]
	}

	matchedUrls := regexp.MustCompile("http.?://.*.(png|jpe?g|gif)").FindStringSubmatch(msg.Message)
	if len(matchedUrls) > 0 {
		virtualBoard = downloadImage(playlist, matchedUrls[0], msg.Inverted, msg.BWThreshold)
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

func downloadImage(playlist *Playlist, imgUrl string, invertImage bool, bwThreshold int) VirtualBoard {
	resp, err := http.Get(imgUrl)
	defer resp.Body.Close()
	m, _, err := image.Decode(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	maxWidth := uint(playlist.PanelInfo.PanelHeight * len(playlist.PanelAddressesLayout[0]))
	maxHeight := uint(playlist.PanelInfo.PanelWidth * len(playlist.PanelAddressesLayout))
	m = resize.Thumbnail(maxWidth, maxHeight, m, resize.Lanczos3)
	bounds := m.Bounds()
	fmt.Printf("%#v \n", bounds)
	var virtualImgBoard VirtualBoard
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		row := fontmap.Row{}
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := m.At(x, y).RGBA()
			lum := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
			pixel := color.Gray{uint8(lum / 256)}

			var flipdotPixelValue bool

			if pixel.Y < uint8(bwThreshold) {
				flipdotPixelValue = !invertImage
			} else {
				flipdotPixelValue = invertImage
			}

			if flipdotPixelValue {
				row = append(row, 1)
			} else {
				row = append(row, 0)
			}
		}
		virtualImgBoard = append(virtualImgBoard, row)
	}
	return virtualImgBoard
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

func startSlackListener(slackToken string, flipboardMsgChn chan FlipBoardDisplayOptions) {
	api := slack.New(slackToken)
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	api.SetDebug(false)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	var oldStopper chan struct{}

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			// close the hold handleSlackMsg
			if oldStopper != nil {
				close(oldStopper)
			}

			stopper := make(chan struct{})
			go handleSlackMsg(ev, rtm, flipboardMsgChn, stopper)
			oldStopper = stopper

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return

		default:
			fmt.Println("Event Received: ")
			fmt.Printf("%#v\n", msg)
		}
	}
}

func handleSlackMsg(slackEvent *slack.MessageEvent, rtm *slack.RTM, flipboardMsgChn chan FlipBoardDisplayOptions, stopper chan struct{}) {
	rawMsg := slackEvent.Msg.Text
	if slackEvent.SubMessage != nil {
		rawMsg = slackEvent.SubMessage.Text
	}

	fmt.Printf("Raw Slack Message: %+v\n", rawMsg)

	messages := splitMessageAndOptions(rawMsg)

	fmt.Printf("%#v \n", messages)

	for _, msg := range messages {
		msg.Message = renderSlackUsernames(msg.Message, rtm)
		msg.Message = cleanupSlackEncodedCharacters(msg.Message)
		msg.Message = renderSlackEmojis(msg.Message, rtm)
		if strings.ToLower(msg.Message) == "help" {
			respondWithHelpMsg(rtm, slackEvent.Msg.Channel)
			return
		}
	}

	for {
		select {
		case <-stopper:
			break // we've received a new message, let's stop looping
		default:
			for _, msg := range messages {
				fmt.Printf("Rendering Message: %+v\n", msg.Message)

				flipboardMsgChn <- msg
				time.Sleep(time.Millisecond * time.Duration(msg.DisplayTime))
			}
		}
	}
}

func splitMessageAndOptions(rawMsg string) ([]FlipBoardDisplayOptions) {
	var messages []FlipBoardDisplayOptions

	playlistRegex := regexp.MustCompile(`^---\n`)
	isPlaylist := playlistRegex.Match([]byte(rawMsg))
	if isPlaylist == false {
		msgAndOptions := regexp.MustCompile(`\s*--(-+)\s*`).Split(rawMsg, -1)
		var m FlipBoardDisplayOptions

		rawOptions := ""
		if len(msgAndOptions) > 1 { // is there options?
			rawOptions = msgAndOptions[1]
		}

		m = unmarshleOptions(rawOptions)
		m.Message = msgAndOptions[0]

		messages = append(messages, m)
	} else {
		rawPlaylist := playlistRegex.Split(rawMsg, -1)[1]

		err := yaml.Unmarshal([]byte(rawPlaylist), &messages)

		if err != nil {
			fmt.Println("Could not unmarshal the yaml")
			fmt.Println(err)
		}
	}

	return messages
}

func cleanupSlackEncodedCharacters(msg string) string {
	// replace slack tokens that are rendered to characters
	msg = strings.Replace(msg, "&lt;", "<", -1)
	msg = strings.Replace(msg, "&gt;", ">", -1)
	return msg
}

func (s *FlipBoardDisplayOptions) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type optsDefaults FlipBoardDisplayOptions

	// todo, make this so that we don't have 2 defaults..
	raw := optsDefaults{
		DisplayTime: 3000,
		Inverted:    false,
		BWThreshold: 140, // magic
		Fill:        "",
		Align:       "center center",
	}

	if err := unmarshal(&raw); err != nil {
		return err
	}

	raw.xAlign, raw.yAlign = getAlignOptions(raw.Align)

	*s = FlipBoardDisplayOptions(raw)
	return nil
}

func unmarshleOptions(rawOptions string) FlipBoardDisplayOptions {
	// reset the options for each Message
	opts := FlipBoardDisplayOptions{
		DisplayTime: 2000,
		Inverted:    false,
		BWThreshold: 140, // magic
		Fill:        "",
		Align:       "center center",
	}

	yaml.Unmarshal([]byte(rawOptions), &opts)

	opts.xAlign, opts.yAlign = getAlignOptions(opts.Align)
	return opts
}

func getAlignOptions(align string) (string, string) {
	alignmentOptions := regexp.MustCompile("( |,)+").Split(align, -1)
	var xAlign, yAlign string
	xAlign = alignmentOptions[0]
	if len(alignmentOptions) > 1 {
		yAlign = alignmentOptions[1]
	}

	return xAlign, yAlign
}

func renderSlackUsernames(msg string, rtm *slack.RTM) string {
	userIds := regexp.MustCompile("<@\\w+>").FindAllString(msg, -1)
	for _, slackFmtMsgUserId := range userIds {
		// in the Message we'll receive something like "<@U123123>", the id is actually "U123123"
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
	msg := `DM me something and I'll try to display that on the board.

You can also supply options for the board by doing:
`

	msg += "```"
	msg += `
Your Message, 🚀, or img_url goes here.
---
align:        # 10 5           // set position of media; horizontally or vertically
align:        # center center  // (left,center,right)  (top,center,bottom)
inverted:     # (true/false) invert the text or image
bwThreshold:  # (0-256) set the threshold value for either "on" or "off"
fill:         # ("", true/false) leave blank for autofill, or select your own fill
`

// we would like to add support for this in the future
//kerning: 0	         // spacing between letters
//font-size: 1           // ??

	msg += "```\n\n"
	rtm.SendMessage(rtm.NewOutgoingMessage(msg, channelId))
}

var slackEmojiLookup map[string]string
var githubEmojiLookup github.EmojiLookup

func renderSlackEmojis(msg string, rtm *slack.RTM) string {
	var err error

	if slackEmojiLookup == nil {
		slackEmojiLookup, err = rtm.GetEmoji()
		if err != nil {
			log.Panicln("Could not get emojis from Slack", err)
			return msg
		}
	}

	emojis := regexp.MustCompile(":\\w+:").FindAllString(msg, -1)
	for _, slackFmtMsgEmoji := range emojis {
		// in the Message we'll receive something like ":smile:", this will actually return 😊
		emojiName := strings.Replace(strings.Replace(slackFmtMsgEmoji, ":", "", 1), ":", "", 1)

		if emojiName != "" {
			emojiImgUrl := slackEmojiLookup[emojiName]

			// follow the aliases for emojis
			for strings.Contains(emojiImgUrl, "alias:") {
				nextEmojiName := strings.Replace(emojiImgUrl, "alias:", "", -1)
				emojiImgUrl = slackEmojiLookup[nextEmojiName]
			}

			if emojiImgUrl == "" {
				emojiImgUrl = githubEmojiLookup[emojiName]
			}

			if emojiImgUrl == "" {
				continue
			}

			msg = strings.Replace(msg, ":"+emojiName+":", emojiImgUrl, -1)
		}
	}

	return msg
}

func findOffSets(options FlipBoardDisplayOptions, virtualBoard VirtualBoard, boardWidth, boardHeight int) (int, int) {
	var xOffSet int

	switch options.xAlign {
	case "left":
		xOffSet = 0
	case "center":
		xOffSet = (boardWidth - len(virtualBoard[0])) / 2
		fmt.Println(xOffSet)
	case "right":
		xOffSet = boardWidth - len(virtualBoard[0])
	default:
		xOffSet, _ = strconv.Atoi(options.xAlign)
	}

	var yOffSet int

	switch options.yAlign {
	case "top":
		// we don't do anything
	case "center":
		yOffSet = (boardHeight - len(virtualBoard)) / 2
	case "bottom":
		yOffSet = boardHeight - len(virtualBoard)
	default:
		yOffSet, _ = strconv.Atoi(options.yAlign)
	}

	return xOffSet, yOffSet
}
