package main

import (
  "github.com/armory/flipdisks/server/internal/health"
  "github.com/armory/flipdisks/server/internal/installations"
  "github.com/gin-contrib/cors"
  "github.com/gin-gonic/gin"
)


func main() {
  router := gin.Default()
  router.Use(cors.Default())

  router.GET("/health", health.Response)

  v1 := router.Group("/v1")
  {
    //v1.GET( "/sites/:siteName/playlists", installations)
    v1.GET("/sites/:siteName/playing", installations.GetPlaying)
    //endpoints for fonts
    v1.GET("/sites/:siteName/fonts")
    //
    v1.POST("/sites/:siteName/fonts/test")
  }

  // listen and serve on 0.0.0.0:8080
  router.Run()
}
