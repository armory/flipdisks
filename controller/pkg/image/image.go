package image

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/armory/flipdisks/pkg/fontmap"
	"github.com/armory/flipdisks/pkg/virtualboard"
	"github.com/nfnt/resize"
	log "github.com/sirupsen/logrus"
)

func ConvertImageUrlToVirtualBoard(maxWidth, maxHeight uint, imgUrl string, invertImage bool, bwThreshold int) *virtualboard.VirtualBoard {
	resp, err := http.Get(imgUrl)
	if err != nil {
		log.Errorf("couldn't download an image %v", err)
		return nil
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		log.Errorf("couldn't decode image %v", err)
		return nil
	}
	defer resp.Body.Close()

	img = resize.Thumbnail(maxWidth, maxHeight, img, resize.Lanczos3)
	bounds := img.Bounds()
	fmt.Printf("%#v \n", bounds)

	return convertImgToVirtualBoard(img, bounds, invertImage, bwThreshold)
}

func convertImgToVirtualBoard(m image.Image, bounds image.Rectangle, invertImage bool, bwThreshold int) *virtualboard.VirtualBoard {
	var board []fontmap.Row
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		row := fontmap.Row{}
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := m.At(x, y).RGBA()
			// to get luminosity, we're going to use magic values from
			// https://stackoverflow.com/questions/596216/formula-to-determine-brightness-of-rgb-color
			lum := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)

			pixel := color.Gray{Y: uint8(lum / 256)} // determine 8 bit gray scale

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
		board = append(board, row)
	}

	v := virtualboard.VirtualBoard(board)
	return &v
}

type FlipboardGif struct {
	Flipboards []*virtualboard.VirtualBoard
	Delay      []time.Duration
}

func convertGifToVirtualBoard(raw []byte, maxWidth, maxHeight uint, invertImage bool, bwThreshold int) (*FlipboardGif, error) {
	flipboardGif := FlipboardGif{
		Flipboards: []*virtualboard.VirtualBoard{},
		Delay:      []time.Duration{},
	}

	g, err := gif.DecodeAll(bytes.NewBuffer(raw))
	if err != nil {
		return &flipboardGif, errors.New("couldn't decode gif: " + err.Error())
	}

	// Create a new RGBA image to hold the incremental frames.
	firstFrame := g.Image[0].Bounds()
	b := image.Rect(0, 0, firstFrame.Dx(), firstFrame.Dy())

	//fmt.Println(b)
	//fmt.Println("number of frames:", len(g.Image))
	//fmt.Println(g.Disposal)
	//fmt.Println("----------------")
	//return &FlipboardGif{}, nil

	// Resize each f.
	lastFrameThatWasntSetToDisposalPrev := 0
	for frameIndex := range g.Image {
		bounds := g.Image[frameIndex].Bounds()
		img := image.NewRGBA(b)
		var newGifFrame *image.Paletted

		//fmt.Println("frameNumber", frameIndex, bounds)

		if frameIndex == 0 {
			draw.Draw(img, bounds, g.Image[frameIndex], bounds.Min, draw.Over)
			lastFrameThatWasntSetToDisposalPrev = frameIndex
		} else {
			frameDisposal := g.Disposal[frameIndex-1] // the last disposal tells us what to do what the base img should be

			switch frameDisposal {
			case byte(0):
				draw.Draw(img, bounds, g.Image[frameIndex], bounds.Min, draw.Over) // just display it
				lastFrameThatWasntSetToDisposalPrev = frameIndex

			case gif.DisposalNone:
				//fmt.Println(frameIndex, "none")

				draw.Draw(img, bounds, g.Image[frameIndex-1], bounds.Min, draw.Over) // get the background
				draw.Draw(img, bounds, g.Image[frameIndex], bounds.Min, draw.Over)   // draw our current frame on top of the background
				lastFrameThatWasntSetToDisposalPrev = frameIndex

			case gif.DisposalBackground:
				//fmt.Println(frameIndex, "background")
				draw.Draw(img, bounds, g.Image[frameIndex], bounds.Min, draw.Over)
				lastFrameThatWasntSetToDisposalPrev = frameIndex

			case gif.DisposalPrevious:
				//	fmt.Println(frameIndex, "prev")
				draw.Draw(img, bounds, g.Image[lastFrameThatWasntSetToDisposalPrev], bounds.Min, draw.Over)
			}
		}

		newGifFrame = imgToPaletted(resizeImage(maxWidth, maxHeight, img))
		vBoard := convertImgToVirtualBoard(newGifFrame, newGifFrame.Bounds(), invertImage, bwThreshold)
		flipboardGif.Flipboards = append(flipboardGif.Flipboards, vBoard)

		// gif time duration is 100th of a second, instead, lets convert it to a time.Duration so it's easier to understand
		flipboardGif.Delay = append(flipboardGif.Delay, time.Duration(g.Delay[frameIndex]/100)*time.Second)

		//fmt.Println("summary:")
		//fmt.Println(g.Image[frameIndex].Bounds())
		fmt.Println(vBoard)

		//return &FlipboardGif{}, nil
		//time.Sleep(time.Millisecond * 500)
	}

	return &flipboardGif, nil
}

func resizeImage(width, height uint, img image.Image) image.Image {
	return resize.Resize(width, height, img, resize.Lanczos3)
}

func imgToPaletted(img image.Image) *image.Paletted {
	b := img.Bounds()
	pm := image.NewPaletted(b, palette.Plan9)
	draw.FloydSteinberg.Draw(pm, b, img, image.ZP)
	return pm
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

	//fmt.Println("asdfasdf")

	return convertGifToVirtualBoard(raw, maxWidth, maxHeight, invertImage, bwThreshold)
}

// GetGifUrl given a message, it'll return an array of gif string urls
func GetGifUrl(url string) []string {
	matched := regexp.MustCompile(`https?://.*\.gif(?:` + `(?:\?(?:\w|\d|&|=|-)+)` + `|` + `(?:\#(?:\w|-)+)` + `)*`).FindStringSubmatch(url)

	// we really don't care about the empty ones
	urls := matched[:0]
	for _, x := range matched {
		if x != "" {
			urls = append(urls, x)
		}
	}

	return urls
}

func GetPlainImageUrl(url string) []string {
	matchImageUrls := regexp.MustCompile(`https?://.*\.(?:png|jpeg|jpg)(?:` + `(?:\?(?:\w|\d|&|=|-)+)` + `|` + `(?:\#(?:\w|-)+)` + `)*`).FindStringSubmatch(url)

	// we really don't care about the empty ones
	urls := matchImageUrls[:0]
	for _, x := range matchImageUrls {
		if x != "" {
			urls = append(urls, x)
		}
	}

	return urls
}
