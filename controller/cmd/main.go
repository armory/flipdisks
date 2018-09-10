package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"github.com/armory/flipdisks/controller/pkg/flipboard"
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

	port := flag.String("p", "/dev/tty.SLAB_USBtoUART", "the serial port, empty string to simulate")
	baud := flag.Int("b", 9600, "baud rate of port")

	var slackToken string
	flag.StringVar(&slackToken, "slack-token", "", "Go get a slack token")
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

	physicalPanels, err := createPanels(playlist, port, baud)
	if err != nil {
		log.Fatal(err)
	}

	// currently we're only supporting uniform panels, oriented the same way
	panelInfo := flipboard.PanelInfo{
		PanelWidth:               28,
		PanelHeight:              7,
		PhysicallyDisplayedWidth: 7,
	}

	panelLayout := [][]flipboard.PanelAddress{
		{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		{10, 11, 12, 13, 14, 15, 16, 17, 18, 19},
	}

	var flipboardOpts []flipboard.Opts
	if countdownDate != "" {
		flipboardOpts = append(flipboardOpts, flipboard.CountdownDate(countdownDate))
	}

	board := flipboard.NewFlipboard(physicalPanels, panelInfo, panelLayout, flipboardOpts...)

	slack := slackbot.NewSlack(slackToken, countdownDate, githubEmojiLookup)

	msgsChan := make(chan options.FlipboardMessageOptions)

	go slack.StartSlackListener(board, msgsChan)

	go flipboard.Play(board)

	time.Sleep(100 * time.Hour)
}


func createPanels(playlist *Playlist, port *string, baud *int) (*[][]panel.Panel, error) {
	var panels [][]panel.Panel

	for y, row := range playlist.PanelAddressesLayout {
		panels = append(panels, []panel.Panel{})

		for _, panelAddress := range row {
			p, err := panel.NewPanel(playlist.PanelInfo.PanelWidth, playlist.PanelInfo.PanelHeight, *port, *baud)
			if err != nil {
				return nil, err
			}

			p.Address = []byte{byte(panelAddress)}

			panels[y] = append(panels[y], *p)
		}
	}
	return &panels, nil
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

