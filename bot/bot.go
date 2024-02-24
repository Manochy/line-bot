package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Manochy/line-bot/models"
)

func Respond(message *models.Message) error {
	msg := strings.ToLower(strings.TrimSpace(message.Text))
	if msg == "c" {
		reply := models.ReplyMessage{
			ReplyToken: message.ReplyToken,
			Messages: []models.MessageResponse{
				{
					Type: "text",
					Text: fmt.Sprintf("Your credit: %.2f", message.Credit),
				},
			},
		}

		if err := SendReply(reply); err != nil {
			return err
		}
	}
	return nil
}

func SendReply(reply models.ReplyMessage) error {
	url := "https://api.line.me/v2/bot/message/reply"

	body, err := json.Marshal(reply)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer YOUR_CHANNEL_ACCESS_TOKEN")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
