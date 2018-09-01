package main

import (
	"github.com/armory/flipdisks/controller/pkg/image"
	"fmt"
)

func main() {
	blah, _ := image.ConvertGifFromURLToVirtualBoard(50, 50, "https://media.giphy.com/media/xUOwGiHZ6NRZfEYYaA/giphy.gif", false, 100)
	for {
		for _, b := range blah.Flipboards {
			s := ""
			for _, a := range *b {
				for _, x := range a {
					if x == 1 {
						s += "⚫️"
					} else {
						s += "⚪️"
					}
				}
				s += "\n"

			}
			fmt.Println(s)
		}
		//time.Sleep(1 * time.Second)
	}
}
