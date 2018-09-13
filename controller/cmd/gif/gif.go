package main

import (
	"github.com/armory/flipdisks/controller/pkg/image"
)

func main() {
	//[2 1 1 2 2 1 2 2 1 1 2 2 1 1 2 2 1 2 2 2 1 2 2 1 1 2 2 1 1 2 2 1 2 2 2 1 2 2 1 1 2 2 2 1 2 1 2 1 2 1 2 2 1 2 1 2 1 2 2 1]
	//url :=  "https://emojis.slackmojis.com/emojis/images/1483822270/1567/drink.gif?1483822270" // drink


	//[2 2 2 2 2 2 2 2 2 1]
	//url := "https://emojis.slackmojis.com/emojis/images/1471119456/981/fast_parrot.gif?1471119456" //parrot	// all 2 and a 1


	//url := "https://i.imgur.com/zb6yJrD.gif" // batman	// all 1s


	url:="https://orig00.deviantart.net/9a18/f/2015/360/3/6/rm_morph_by_bitsandpieces12-d9llus2.gif" // rick/morty
	//url:= "https://thumbs.gfycat.com/ShoddyLeftKiwi-size_restricted.gif" // 123    all zeros
	//url:="https://media.giphy.com/media/emP6pgjuDMQOA/source.gif" // all ones
	//url := "https://media.giphy.com/media/AaVVwrwfIlTa0/giphy.gif"
	//url := "https://media.giphy.com/media/xUOwGiHZ6NRZfEYYaA/giphy.gif" // fish, all 1

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
		return
	}
}
