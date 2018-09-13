package main

import (
	"github.com/armory/flipdisks/controller/pkg/image"
)

func main() {
	url :=  "https://emojis.slackmojis.com/emojis/images/1483822270/1567/drink.gif?1483822270" // drink
	//url := "https://emojis.slackmojis.com/emojis/images/1471119456/981/fast_parrot.gif?1471119456" //parrot

	//blah, _ := image.ConvertGifFromURLToVirtualBoard( url,50, 50,  false, 90)
	for {
		image.ConvertGifFromURLToVirtualBoard( url,50, 50,  false, 90)
		//for _, b := range blah.Flipboards {
		//	s := ""
		//	for _, a := range *b {
		//		for _, x := range a {
		//			if x == 1 {
		//				s += "⚫️"
		//			} else {
		//				s += "⚪️"
		//			}
		//		}
		//		s += "\n"
		//
		//	}
			//fmt.Println(s)
			//time.Sleep(time.Duration(1) * time.Second)
		//}
	}
}
