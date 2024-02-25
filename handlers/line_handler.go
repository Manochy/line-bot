package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

const accessToken = "iGraIEph4ShvuNEYaLVaPLnwHIKcQvzWM3T8mdw8ir4CwCREUMcp6XdEk/L88evFEPV9qnAWGgGXqE5GNW720bi137Ydc5hX6l1bL1K3N9nEp08KBsK8CrDzEZkbmxRTs9Ns6lUG5lnKiFKecxd0fAdB04t89/1O/w1cDnyilFU="
const internalErrTxt = "Internal server error"

type Request struct {
	Events []struct {
		ReplyToken string `json:"replyToken"`
		Type       string `json:"type"`
		Message    struct {
			ID   string `json:"id"`
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"message"`
		Source struct {
			UserID string `json:"userId"`
		} `json:"source"`
	} `json:"events"`
}

type Response struct {
	ReplyToken string    `json:"replyToken"`
	Messages   []Message `json:"messages"`
}

type Message struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func LineBotHandler(c *gin.Context) {
	var req Request
	if err := c.BindJSON(&req); err != nil {
		log.Println("Error decoding request body:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": internalErrTxt})
		return
	}

	db, err := sql.Open("mysql", "root:Evild0ergu@tcp(localhost:3306)/pokerth")
	if err != nil {
		log.Println("Error connecting to database:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": internalErrTxt})
		return
	}
	defer db.Close()

	for _, event := range req.Events {
		if event.Type == "message" && event.Message.Type == "text" {

			// COMMAND "test"
			if event.Message.Text == "test" {
				println(" >>>>>>>>>>>> User ID is " + event.Source.UserID)
				replyMessage := Response{
					ReplyToken: event.ReplyToken,
					Messages: []Message{
						{
							Type: "text",
							Text: "Your LINE ID is: " + event.Source.UserID,
						},
					},
				}

				if err := replyToLine(replyMessage, accessToken); err != nil {
					log.Println("Error sending reply message to LINE:", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": internalErrTxt})
					return
				}
			}

			if event.Message.Text == "c" {
				var userCredit float64
				err := db.QueryRow("SELECT member_credit FROM member WHERE member_line_Id = ?", event.Source.UserID).Scan(&userCredit)
				if err != nil {
					log.Println("Error retrieving user credit from database:", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": internalErrTxt})
					return
				}

				replyMessage := Response{
					ReplyToken: event.ReplyToken,
					Messages: []Message{
						{
							Type: "text",
							Text: "Your credit is: " + strconv.FormatFloat(userCredit, 'f', 2, 64),
						},
					},
				}

				if err := replyToLine(replyMessage, accessToken); err != nil {
					log.Println("Error sending reply message to LINE:", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": internalErrTxt})
					return
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

func replyToLine(replyMessage Response, accessToken string) error {
	// Convert replyMessage to JSON
	replyJSON, err := json.Marshal(replyMessage)
	if err != nil {
		return err
	}

	log.Println("Reply JSON:", string(replyJSON))

	// Create HTTP request
	req, err := http.NewRequest("POST", "https://api.line.me/v2/bot/message/reply", bytes.NewBuffer(replyJSON))
	if err != nil {
		return err
	}

	// Add authorization header with access token
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil

}
