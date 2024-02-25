package main

import (
	"log"
	"net/http"

	"github.com/Manochy/line-bot/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/webhook", handlers.WebhookHandler)

	log.Println("Starting server on port 8080...")
	if err := http.ListenAndServe("0.0.0.0:8080", r); err != nil {
		log.Fatal("Server error:", err)
	}
}
