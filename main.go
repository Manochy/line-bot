// main.go

package main

import (
	"log"
	"os"

	"github.com/Manochy/line-bot/handlers"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client

func main() {
	var err error
	bot, err = linebot.New(
		"43d644a6edbf6bbcf7f30b29672ff6be",
		"3sPiJyEc0px9rXJiv80h5HAvoqTF+RF5pGqfIcsh/5vEZ8+x+b2E064u8Yt4f5o2QwFhRit8G9Rn034AKMG4Z6Bxiur0qX7w9mCuasUlYerYW9H1D4b8sFLjQIxyN9cPW5VzbdxS/3FuYJKONzM+ewdB04t89/1O/w1cDnyilFU=",
	)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.POST("/callback", handlers.HandleCallback(bot))

	port := os.Getenv("PORT")
	if port == "" {
		port = "12024"
	}
	log.Printf("Starting server on port %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
