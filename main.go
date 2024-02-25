package main

import (
	"log"
	"os"

	"github.com/Manochy/line-bot/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/callback", handlers.LineBotHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server listening on port %s", port)
	r.Run(":" + port)
}
