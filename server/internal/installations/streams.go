package installations

import (
	"github.com/gin-gonic/gin"
	"log"
)


type Row []int
type Board []Row
type Frame map[string]Board


type streamResponse struct {
	Installation string     `json:"installation"`
	FrameRate    int        `json:"frameRate"`
	Looping      bool       `json:"looping"`
	Layout			 [][]string `json:"layout"`
	Frames 			 []Frame    `json:"frames"`
}


func GetStream(c *gin.Context) {
	installationLocation := c.Param("installationLocation")

	log.Println("sending payload stub for: ", installationLocation)

	stream := streamResponse{
		Installation: installationLocation,
		FrameRate:    1000,
		Looping:      true,
		Layout:				[][]string {
	    []string{"board1", "board2", "board3"},
	    []string{"board4", "board5", "board6"},
	    []string{"board7", "board8", "board9"},
	    []string{"board7", "board8", "board9"},
  	},
		Frames: []Frame{
			Frame{
				"board1": Board{
					Row{1, 1, 1, 1, 1, 1},
					Row{1,1,1,1,1,1,1},
				},
			},
		},
	}

	c.JSON(200, stream)
}
