package installations

import (
	"github.com/gin-gonic/gin"
	"log"
)


type streamResponse struct {
	Installation string `json:"installation"`
	FrameRate    int    `json:"frameRate"`
	Looping      bool   `json:"looping"`
}


func GetStream(c *gin.Context) {
	installationLocation := c.Param("installationLocation")

	log.Println("sending payload stub for: ", installationLocation)

	stream := streamResponse{
		Installation: installationLocation,
		FrameRate:    1000,
		Looping:      true,
	}

	c.JSON(200, stream)
}
