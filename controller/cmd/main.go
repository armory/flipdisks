package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/armory/flipdisks/controller/pkg/flipboard"
	"github.com/armory/flipdisks/controller/pkg/github"
	"github.com/armory/flipdisks/controller/pkg/slackbot"
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

	// currently we're only supporting uniform panels, oriented the same way
	panelInfo := flipboard.PanelInfo{
		PanelWidth:               28,
		PanelHeight:              7,
		PhysicallyDisplayedWidth: 7,
		Port: *port,
		Baud: *baud,
	}

	panelLayout := [][]flipboard.PanelAddress{
		{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		{10, 11, 12, 13, 14, 15, 16, 17, 18, 19},
	}

	var flipboardOpts []flipboard.Opts
	if countdownDate != "" {
		flipboardOpts = append(flipboardOpts, flipboard.CountdownDate(countdownDate))
	}

	board, err := flipboard.NewFlipboard(panelInfo, panelLayout, flipboardOpts...)
	if err != nil {
		log.Fatal("couldn't create flipboard: " + err.Error())
	}

	slack := slackbot.NewSlack(slackToken, countdownDate, githubEmojiLookup)

	go slack.StartSlackListener(board)

	go flipboard.Play(board)

	// we're actually going to just block forever so the program stays alive
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
