package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"

	"github.com/armory/flipdisks/controller/pkg/flipboard"
	"github.com/armory/flipdisks/controller/pkg/github"
	myimage "github.com/armory/flipdisks/controller/pkg/image"
	"github.com/armory/flipdisks/controller/pkg/options"
	"github.com/nfnt/resize"
	log "github.com/sirupsen/logrus"

	"gocv.io/x/gocv"
)

func main() {
	log.Print("Starting flipdisk controller")

	port := flag.String("p", "/dev/tty.SLAB_USBtoUART", "the serial port, empty string to simulate")
	baud := flag.Int("b", 9600, "baud rate of port")

	var slackToken string
	flag.StringVar(&slackToken, "slack-token", "", "Go get a slack token")

	var githubToken string
	flag.StringVar(&githubToken, "github-token", "", "Go get a github token")

	var countdownDate string
	flag.StringVar(&countdownDate, "countdown", "", fmt.Sprintf("Specify the countdown date in YYYY-MM-DD format"))
	flag.Parse()

	g, err := github.New(github.Token(githubToken))
	if err != nil {
		log.Error("Could not create githubClient, hopefully everything will work!")
	}
	_, err = g.GetEmojis()
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

	//start
	deviceID := 0
	webcam, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		fmt.Printf("Error opening video capture device: %v\n", deviceID)
		return
	}
	webcam.Set(gocv.VideoCaptureFPS, 10)
	webcam.Set(gocv.VideoCaptureFrameWidth, 50)
	webcam.Set(gocv.VideoCaptureFrameHeight, 50)
	defer webcam.Close()

	window := gocv.NewWindow("Motion Window")
	defer window.Close()

	img := gocv.NewMat()
	defer img.Close()

	imgDelta := gocv.NewMat()
	defer imgDelta.Close()

	imgThresh := gocv.NewMat()
	defer imgThresh.Close()

	mog2 := gocv.NewBackgroundSubtractorMOG2()
	defer mog2.Close()
	//

	fmt.Printf("Start reading device: %v\n", deviceID)
	opp := options.GetDefaultOptions()
	opp.DisplayTime = 0
	for {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("Device closed: %v\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}

		mog2.Apply(img, &imgDelta)
		gocv.Threshold(imgDelta, &imgThresh, 10, 255, gocv.ThresholdBinary)
		kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(3, 3))
		defer kernel.Close()
		gocv.Erode(imgThresh, &imgThresh, kernel)

		q, err := gocv.IMEncode(gocv.PNGFileExt, img)
		if err != nil {
			panic(err)
		}
		qq, _, err := image.Decode(bytes.NewReader(q))
		if err != nil {
			panic(err)
		}
		qq = resize.Thumbnail(50, 50, qq, resize.Lanczos3)
		zz := myimage.ConvertImgToVirtualBoard(qq, qq.Bounds(), false, 90)
		flipboard.DisplayVirtualBoardToPhysicalBoard(&opp, zz, board)
	}
}
