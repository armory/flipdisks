package image

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"regexp"

	"github.com/armory/flipdisks/controller/pkg/fontmap"
	"github.com/nfnt/resize"
	log "github.com/sirupsen/logrus"
	"time"
	"image/gif"
	"errors"
	"io/ioutil"
	"github.com/discordapp/lilliput"
	"bytes"
)

func Download(maxWidth, maxHeight uint, imgUrl string, invertImage bool, bwThreshold int) []fontmap.Row {
	resp, err := http.Get(imgUrl)
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		log.Errorf("couldn't download an image %v", err)
		return nil
	}
	defer resp.Body.Close()

	img = resize.Thumbnail(maxWidth, maxHeight, img, resize.Lanczos3)
	bounds := img.Bounds()
	fmt.Printf("%#v \n", bounds)

	return convertImgToVirtualBoard(img, bounds, invertImage, bwThreshold)
}

func convertImgToVirtualBoard(m image.Image, bounds image.Rectangle, invertImage bool, bwThreshold int) []fontmap.Row {
	var virtualImgBoard []fontmap.Row
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		row := fontmap.Row{}
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// use magic values from
			// https://stackoverflow.com/questions/596216/formula-to-determine-brightness-of-rgb-color
			r, g, b, _ := m.At(x, y).RGBA()
			lum := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)

			pixel := color.Gray{Y: uint8(lum / (2 ^ 8))} // determine 8 bit gray scale

			flipdotPixelValue := invertImage
			if pixel.Y < uint8(bwThreshold) {
				flipdotPixelValue = !invertImage
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

type FlipboardGif struct {
	Flipboards []*[]fontmap.Row
	Delay      []time.Duration
}

func convertGifToVirtualBoard(raw []byte, maxWidth, maxHeight uint, invertImage bool, bwThreshold int) (*FlipboardGif, error) {
	flipboardGif := FlipboardGif{
		Flipboards: []*[]fontmap.Row{},
		Delay:      []time.Duration{},
	}

	decoder, err := lilliput.NewDecoder(raw)
	defer decoder.Close()
	if err != nil {
		fmt.Println(err)
		return &flipboardGif, errors.New("couldn't decode raw gif data" + err.Error())
	}

	//resizing image
	ops := lilliput.NewImageOps(8192) // magic size
	defer ops.Close()

	opts := &lilliput.ImageOptions{
		FileType:             ".gif",
		Width:                int(maxWidth),
		Height:               int(maxHeight),
		ResizeMethod:         lilliput.ImageOpsFit,
		NormalizeOrientation: true,
	}

	outputImg := make([]byte, 50*1024*1024) // magic size
	outputImg, err = ops.Transform(decoder, opts, outputImg)
	if err != nil {
		return &flipboardGif, errors.New("couldn't resize image: " + err.Error())
	}
	// pass downloaded iwmage to gif.DecodeAll, return *GIF
	gif, err := gif.DecodeAll(bytes.NewBuffer(outputImg))

	if err != nil {
		return &flipboardGif, errors.New("couldn't decode gif: " + err.Error())
	}

	// for each frame in a gif
	for i, frame := range gif.Image {
		//    resize each frame
		fmt.Println(frame.Bounds())
		//    call convertImgToVirtualBoard() to return a flipboard
		vBoard := convertImgToVirtualBoard(frame, frame.Bounds(), invertImage, bwThreshold)
		//    append to flipboardGif
		flipboardGif.Flipboards = append(flipboardGif.Flipboards, &vBoard)
		flipboardGif.Delay = append(flipboardGif.Delay, time.Duration(gif.Delay[i]/100)*time.Second)

	}

	return &flipboardGif, nil
}

func ConvertGifFromURLToVirtualBoard(gifUrl string, maxWidth, maxHeight uint, invertImage bool, bwThreshold int) (*FlipboardGif, error) {
	// download the image http.Get
	r, err := http.Get(gifUrl)
	if err != nil {
		return &FlipboardGif{}, errors.New("couldn't download gif: " + err.Error())
	}

	raw, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return &FlipboardGif{}, errors.New("couldn't get raw gif data: " + err.Error())
	}

	return convertGifToVirtualBoard(raw, maxWidth, maxHeight, invertImage, bwThreshold)
}

func IsGifUrl(url string) bool {
	matchGifUrls := regexp.MustCompile(`^http.?://.*\.(gif)`).FindStringSubmatch(url)
	if len(matchGifUrls) > 0 {
		return true
	}
	return false
}


func IsPlainImageUrl(url string) bool {
	matchImageUrls := regexp.MustCompile(`^http.?://.*\.(png|jpe?g)`).FindStringSubmatch(url)
	if len(matchImageUrls) > 0 {
		return true
	}
	return false
}

