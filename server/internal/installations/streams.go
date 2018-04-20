package installations

import (
	"github.com/gin-gonic/gin"
	"log"
)

type Row []int
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

func GetPlaying(c *gin.Context) {
	siteName := c.Param("siteName")

	log.Println("sending payload stub for: ", siteName)

	video := FlipdiskVideo{
		Installation: siteName,
		FPS:          1,
		Looping:      true,
		Layout: [][]string{
			[]string{"board1", "board2", "board3"},
			[]string{"board4", "board5", "board6"},
			[]string{"board7", "board8", "board9"},
			[]string{"board7", "board8", "board9"},
		},
		Frames: []Frame{
			Frame{
				"board1": Board{
					Row{1, 1, 1, 1, 1, 1, 1},
					Row{1, 1, 1, 1, 1, 1, 1},
				},
				"board2": Board{
					Row{1, 1, 1, 1, 1, 1, 1},
					Row{1, 1, 1, 1, 1, 1, 1},
				},
				"board3": Board{
					Row{1, 1, 1, 1, 1, 1, 1},
					Row{1, 1, 1, 1, 1, 1, 1},
				},
				"board4": Board{
					Row{1, 1, 1, 1, 1, 1, 1},
					Row{1, 1, 1, 1, 1, 1, 1},
				},
				"board5": Board{
					Row{1, 1, 1, 1, 1, 1, 1},
					Row{1, 1, 1, 1, 1, 1, 1},
				},
				"board6": Board{
					Row{1, 1, 1, 1, 1, 1, 1},
					Row{1, 1, 1, 1, 1, 1, 1},
				},
				"board7": Board{
					Row{1, 1, 1, 1, 1, 1, 1},
					Row{1, 1, 1, 1, 1, 1, 1},
				},
				"board8": Board{
					Row{1, 1, 1, 1, 1, 1, 1},
					Row{1, 1, 1, 1, 1, 1, 1},
				},
				"board9": Board{
					Row{1, 1, 1, 1, 1, 1, 1},
					Row{1, 1, 1, 1, 1, 1, 1},
				},
			},
			Frame{
				"board1": Board{
					Row{2, 2, 2, 2, 2, 2, 2},
					Row{2, 2, 2, 2, 2, 2, 2},
				},
				"board2": Board{
					Row{2, 2, 2, 2, 2, 2, 2},
					Row{2, 2, 2, 2, 2, 2, 2},
				},
				"board3": Board{
					Row{2, 2, 2, 2, 2, 2, 2},
					Row{2, 2, 2, 2, 2, 2, 2},
				},
				"board4": Board{
					Row{2, 2, 2, 2, 2, 2, 2},
					Row{2, 2, 2, 2, 2, 2, 2},
				},
				"board5": Board{
					Row{2, 2, 2, 2, 2, 2, 2},
					Row{2, 2, 2, 2, 2, 2, 2},
				},
				"board6": Board{
					Row{2, 2, 2, 2, 2, 2, 2},
					Row{2, 2, 2, 2, 2, 2, 2},
				},
				"board7": Board{
					Row{2, 2, 2, 2, 2, 2, 2},
					Row{2, 2, 2, 2, 2, 2, 2},
				},
				"board8": Board{
					Row{2, 2, 2, 2, 2, 2, 2},
					Row{2, 2, 2, 2, 2, 2, 2},
				},
				"board9": Board{
					Row{2, 2, 2, 2, 2, 2, 2},
					Row{2, 2, 2, 2, 2, 2, 2},
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
