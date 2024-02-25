package main

import (
	"github.com/Manochy/line-bot/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/callback", handlers.LineBotHandler)
	r.Run(":12024")
}
