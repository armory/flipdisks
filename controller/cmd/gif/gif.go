package main

import (
	"github.com/armory/flipdisks/controller/pkg/image"
	"fmt"
	"time"
)

func main() {
	blah, _ := image.ConvertGifFromURLToVirtualBoard("https://emojis.slackmojis.com/emojis/images/1471119456/981/fast_parrot.gif?1471119456",50, 50,  false, 90)
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
			time.Sleep(time.Duration(1) * time.Second)
		}
	}
}
