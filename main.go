package main

import (
	"M_Pro/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.PUT("/change/:date/:prayer", handlers.ChangeCH)
	r.GET("/dates/:date", handlers.GetHistory)

	r.Run(":1902")
}
