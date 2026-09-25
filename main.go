package main

import (
	"M_Pro/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/change", handlers.PostCh)
	r.GET("/dates", handlers.GetHistory)

	r.Run(":1712")
}
