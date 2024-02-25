package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// Request represents the JSON structure of the LINE webhook event
type Request struct {
	Events []struct {
		ReplyToken string `json:"replyToken"`
		Type       string `json:"type"`
		Message    struct {
			ID   string `json:"id"`
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"message"`
	} `json:"events"`
}

// Response represents the JSON structure of the LINE reply message
type Response struct {
	ReplyToken string `json:"replyToken"`
	Messages   []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"messages"`
}

// LineBotHandler handles LINE webhook events
func LineBotHandler(c *gin.Context) {
	// Parse request body
	var req Request
	const inServerErrTxt = "Internal server error"
	if err := c.BindJSON(&req); err != nil {
		log.Println("Error decoding request body:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": inServerErrTxt})
		return
	}

	// Connect to MySQL database
	db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/database")
	if err != nil {
		log.Println("Error connecting to database:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": inServerErrTxt})
		return
	}
	defer db.Close()

	// Process each event
	for _, event := range req.Events {
		if event.Type == "message" {
			// Retrieve member credit from MySQL (replace "your_member_id" with actual member ID)
			memberCredit, err := getMemberCredit(db, "your_member_id")
			if err != nil {
				log.Println("Error retrieving member credit from database:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": inServerErrTxt})
				return
			}

			// Construct the reply message
			replyMessage := Response{
				ReplyToken: event.ReplyToken,
				Messages: []struct {
					Type string `json:"type"`
					Text string `json:"text"`
				}{
					{
						Type: "text",
						Text: event.Message.Text + " From Bot, Your Credit: " + strconv.FormatFloat(memberCredit, 'f', 2, 64),
					},
				},
			}

			// Send the reply message to LINE
			if err := replyToLine(replyMessage); err != nil {
				log.Println("Error sending reply message to LINE:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": inServerErrTxt})
				return
			}
		}
	}

	// Respond with status OK
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

func getMemberCredit(db *sql.DB, memberID string) (float64, error) {
	var memberCredit float64
	err := db.QueryRow("SELECT member_credit FROM member WHERE member_id = ?", memberID).Scan(&memberCredit)
	if err != nil {
		return 0, err
	}
	return memberCredit, nil
}

func replyToLine(replyMessage Response) error {
	// Marshal the reply message to JSON
	replyJSON, err := json.Marshal(replyMessage)
	if err != nil {
		return err
	}

	// Send the reply message to LINE
	_, err = http.Post("https://api.line.me/v2/bot/message/reply", "application/json", bytes.NewBuffer(replyJSON))
	if err != nil {
		return err
	}

	return nil
}
