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
)

func Download(maxWidth, maxHeight uint, imgUrl string, invertImage bool, bwThreshold int) []fontmap.Row {
	resp, err := http.Get(imgUrl)
	m, _, err := image.Decode(resp.Body)
	if err != nil {
		log.Errorf("couldn't download an image %v", err)
		return nil
	}
	defer resp.Body.Close()

	m = resize.Thumbnail(maxWidth, maxHeight, m, resize.Lanczos3)
	bounds := m.Bounds()
	fmt.Printf("%#v \n", bounds)

	return convert(m, bounds, invertImage, bwThreshold)
}

func convert(m image.Image, bounds image.Rectangle, invertImage bool, bwThreshold int) []fontmap.Row {
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

