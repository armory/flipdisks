package installations

import (
	"github.com/gin-gonic/gin"
	"log"
)

type Row []string
type Board []Row
type Frame map[string]Board

type FlipdiskVideo struct {
	Installation string     `json:"installation"`
	FPS          int        `json:"fps"`
	Looping      bool       `json:"looping"`
	Layout       [][]string `json:"layout"`
	Frames       []Frame    `json:"frames"`
}

type Playlist struct {
	Name    string          `json:"name"`
	Videos  []FlipdiskVideo `json:"videos"`
	Looping bool            `json:"looping"`
}


// mikronf_5x5 -  https://fontstruct.com/fontstructions/show/74288/mikronf_5x5
// pixelspace_5x5 - https://fontstruct.com/fontstructions/show/42130/
// 5x5 - https://fontstruct.com/fontstructions/show/1916/swos
// ti84_small -  https://fontstruct.com/fontstructions/show/424974/ti84_small
// swos - 5x5 https://fontstruct.com/fontstructions/show/1916/swos


func GetPlaying(c *gin.Context) {
	siteName := c.Param("siteName")

	log.Println("sending payload stub for: ", siteName)

	video := FlipdiskVideo{
		Installation: siteName,
		FPS:          1,
		Looping:      true,
		Layout: [][]string{
			[]string{"1", "2", "3", "4", "5",},
		},
		Frames: []Frame{
			Frame{
				"1": Board{
					Row{"⚫️", "⚪️", "⚪️", "⚫️", "⚪️",},
					Row{"⚫️", "⚪️", "⚪️", "⚫️", "⚪️",},
					Row{"⚫️", "⚪️", "⚪️", "⚫️", "⚪️",},
					Row{"⚫️", "⚫️", "⚫️", "⚫️", "⚪️",},
					Row{"⚫️", "⚪️", "⚪️", "⚫️", "⚪️",},
					Row{"⚫️", "⚪️", "⚪️", "⚫️", "⚪️",},
				},
				"2": Board{
					Row{"⚪️", "⚪️", "⚪️", "⚪️", "⚪️",},
					Row{"⚫️", "⚫️", "⚫️", "⚫️", "⚪️",},
					Row{"⚫️", "⚪️", "⚪️", "⚪️", "⚪️",},
					Row{"⚫️", "⚫️", "⚫️", "⚫️", "⚪️",},
					Row{"⚫️", "⚪️", "⚪️", "⚪️", "⚪️",},
					Row{"⚫️", "⚫️", "⚫️", "⚫️", "⚪️",},
				},
				"3": Board{
					Row{"⚫️", "⚪️", "⚪️", "⚪️", "⚪️",},
					Row{"⚫️", "⚪️", "⚪️", "⚪️", "⚪️",},
					Row{"⚫️", "⚪️", "⚪️", "⚪️", "⚪️",},
					Row{"⚫️", "⚪️", "⚪️", "⚪️", "⚪️",},
					Row{"⚫️", "⚪️", "⚪️", "⚪️", "⚪️",},
					Row{"⚫️", "⚫️", "⚫️", "⚪️", "⚪️",},
				},
				"4": Board{
					Row{"⚫️", "⚪️", "⚪️", "⚪️", "⚪️",},
					Row{"⚫️", "⚪️", "⚪️", "⚪️", "⚪️",},
					Row{"⚫️", "⚪️", "⚪️", "⚪️", "⚪️",},
					Row{"⚫️", "⚪️", "⚪️", "⚪️", "⚪️",},
					Row{"⚫️", "⚪️", "⚪️", "⚪️", "⚪️",},
					Row{"⚫️", "⚫️", "⚫️", "⚪️", "⚪️",},
				},
				"5": Board{
					Row{"⚪️", "⚫️", "⚫️", "⚪️", "⚪️",},
					Row{"⚫️", "⚪️", "⚪️", "⚫️", "⚪️",},
					Row{"⚫️", "⚪️", "⚪️", "⚫️", "⚪️",},
					Row{"⚫️", "⚪️", "⚪️", "⚫️", "⚪️",},
					Row{"⚫️", "⚪️", "⚪️", "⚫️", "⚪️",},
					Row{"⚪️", "⚫️", "⚫️", "⚪️", "⚪️",},
				},
			},
			Frame{
				// w
				"1": Board{
					Row{"⚫️", "⚪️", "⚫️",},
					Row{"⚫️", "⚪️", "⚫️",},
					Row{"⚫️", "⚪️", "⚫️",},
					Row{"⚫️", "⚪️", "⚫️",},
					Row{"⚪️", "⚫️", "⚫️",},
				},

				// o
				"2": Board{
					Row{"⚪️", "⚫️", "⚪️",},
					Row{"⚫️", "⚪️", "⚫️",},
					Row{"⚫️", "⚪️", "⚫️",},
					Row{"⚫️", "⚪️", "⚫️",},
					Row{"⚪️", "⚫️", "⚪️",},
				},

				// r
				"3": Board{
					Row{"⚫️", "⚫️", "⚫️", "⚫️", "⚪️",},
					Row{"⚫️", "⚪️", "⚪️", "⚪️", "⚫️",},
					Row{"⚫️", "⚫️", "⚫️", "⚫️", "⚪️",},
					Row{"⚫️", "⚪️", "⚪️", "⚪️", "⚫️",},
					Row{"⚫️", "⚪️", "⚪️", "⚪️", "⚫️",},
				},

				// l
				"4": Board{
					Row{"⚫️", "⚪️", "⚪️", "⚪️", "⚪️",},
					Row{"⚫️", "⚪️", "⚪️", "⚪️", "⚪️",},
					Row{"⚫️", "⚪️", "⚪️", "⚪️", "⚪️",},
					Row{"⚫️", "⚪️", "⚪️", "⚪️", "⚪️",},
					Row{"⚫️", "⚫️", "⚫️", "⚪️", "⚪️",},
				},

				// d
				"5": Board{
					Row{"⚫️", "⚫️", "⚪️", "⚪️", "⚪️",},
					Row{"⚫️", "⚪️", "⚫️", "⚪️", "⚪️",},
					Row{"⚫️", "⚪️", "⚫️", "⚪️", "⚪️",},
					Row{"⚫️", "⚪️", "⚫️", "⚪️", "⚪️",},
					Row{"⚫️", "⚫️", "⚪️", "⚪️", "⚪️",},
				},
			},
		},
	}

	playlist := Playlist{
		Name: "now showing",
		Videos: []FlipdiskVideo{
			video,
		},
	}

	c.JSON(200, playlist)
}
