package flipimage

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"

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

	m = resize.Thumbnail(20, 20, m, resize.Lanczos3)
	bounds := m.Bounds()
	fmt.Printf("%#v \n", bounds)

	return convert(m, bounds, invertImage, bwThreshold)
}

func convert(m image.Image, bounds image.Rectangle, invertImage bool, bwThreshold int) []fontmap.Row {
	var virtualImgBoard []fontmap.Row
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
