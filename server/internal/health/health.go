package health

import (
	"github.com/gin-gonic/gin"
	"time"
)

func Response(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
		"ts": time.Now(),
	})
}
