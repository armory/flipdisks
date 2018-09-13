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

	"github.com/armory/flipdisks/controller/pkg/fontmap"
	"github.com/armory/flipdisks/controller/pkg/virtualboard"
	"github.com/nfnt/resize"
	log "github.com/sirupsen/logrus"
)

func ConvertImageUrlToVirtualBoard(maxWidth, maxHeight uint, imgUrl string, invertImage bool, bwThreshold int) *virtualboard.VirtualBoard {
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
	img := image.NewRGBA(b)

	// Resize each f.
	lastIndexToHaveNoDisposalPrev := 0
	for frameIndex := range g.Image {
		frameDisposal := g.Disposal[frameIndex]

		switch frameDisposal {
		case gif.DisposalPrevious:
			fmt.Println(frameIndex, "prev")
			g.Image[frameIndex] = g.Image[lastIndexToHaveNoDisposalPrev]
			g.Image[frameIndex] = imgToPaletted(resizeImage(maxWidth, maxHeight, g.Image[frameIndex]))
		case gif.DisposalNone:
			fmt.Println(frameIndex, "none")
			lastIndexToHaveNoDisposalPrev = frameIndex
			g.Image[frameIndex+1] = imgToPaletted(resizeImage(maxWidth, maxHeight, g.Image[frameIndex]))
		case gif.DisposalBackground:
			fmt.Println(frameIndex, "background")
			lastIndexToHaveNoDisposalPrev = frameIndex
			if frameIndex > 0 {
				g.Image[frameIndex] = imgToPaletted(g.Image[frameIndex-1])
			} else {
				g.Image[frameIndex] = imgToPaletted(resizeImage(maxWidth, maxHeight, g.Image[frameIndex]))
			}
		}

		bounds := g.Image[frameIndex].Bounds()
		draw.Draw(img, bounds, g.Image[frameIndex], bounds.Min, draw.Over)

		vBoard := convertImgToVirtualBoard(g.Image[frameIndex], g.Image[frameIndex].Bounds(), invertImage, bwThreshold)
		flipboardGif.Flipboards = append(flipboardGif.Flipboards, vBoard)

		// gif time duration is 100th of a second, instead, lets convert it to a time.Duration so it's easier to understand
		flipboardGif.Delay = append(flipboardGif.Delay, time.Duration(g.Delay[frameIndex]/100)*time.Second)

		fmt.Println(g.Image[frameIndex].Bounds())
		fmt.Println(vBoard)

		time.Sleep(time.Millisecond * 500)
	}

	return &flipboardGif, nil
}

func resizeImage(width, height uint, img image.Image) image.Image {
	return resize.Resize(width, height, img, resize.NearestNeighbor)
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

	return convertGifToVirtualBoard(raw, maxWidth, maxHeight, invertImage, bwThreshold)
}

// GetGifUrl given a message, it'll return an array of gif string urls
func GetGifUrl(url string) []string {
	matched := regexp.MustCompile(`https?://.*\.gif(?:(\\?)\S+)?(?:#\S+)?`).FindStringSubmatch(url)

	// we really don't care about the empty ones
	urls := matched[:0]
	for _, x := range matched {
		if x != "" {
			urls = append(urls, x)
		}
	}

	return urls
}

func IsPlainImageUrl(url string) bool {
	matchImageUrls := regexp.MustCompile(`^http.?://.*\.(?:png|jpe?g)(?:(\\?)\S+)?(?:#\S+)?`).FindStringSubmatch(url)
	if len(matchImageUrls) > 0 {
		return true
	}
	return false
}
