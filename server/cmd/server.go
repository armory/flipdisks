package main

import (
  "github.com/armory/flipdisks/server/internal/health"
  "github.com/armory/flipdisks/server/internal/installations"
  "github.com/gin-gonic/gin"
)


func main() {
  r := gin.Default()

  r.GET("/health", health.Response)

  v1 := r.Group("/v1")
  {
    v1.GET("/installations/:installationLocation/view", installations.GetStream)
  }

  // listen and serve on 0.0.0.0:8080
  r.Run()
}
