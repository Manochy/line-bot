package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Manochy/line-bot/bot"
	"github.com/Manochy/line-bot/models"
)

func WebhookHandler(c *gin.Context) {
	var message models.Message
	if err := c.BindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	response := bot.Respond(&message)
	c.JSON(http.StatusOK, gin.H{"response": response})
}
