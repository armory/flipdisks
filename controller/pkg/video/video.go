package video

import (
	"bytes"
	"fmt"
	"image"

	"github.com/armory/flipdisks/controller/pkg/flipboard"
	myimage "github.com/armory/flipdisks/controller/pkg/image"
	"github.com/armory/flipdisks/controller/pkg/options"
	"github.com/nfnt/resize"
	"gocv.io/x/gocv"
)

type video struct {
	stop chan bool
}

var instance *video

func init() {
	instance = &video{
		stop: make(chan bool, 1),
	}
}

func GetVideo() *video {
	return instance
}

func (v *video) Stop() {
	close(v.stop)
}

func (v *video) Start(board *flipboard.Flipboard) {
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

loop:
	for {
		select {
		case <-v.stop:
			fmt.Println("stopping")
			break loop
		}

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
