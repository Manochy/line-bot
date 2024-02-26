package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Manochy/line-bot/models"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
)

func HandleCallback(bot *linebot.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		events, err := bot.ParseRequest(c.Request)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				c.Status(http.StatusBadRequest)
				return
			}
			c.Status(http.StatusInternalServerError)
			return
		}

		for _, event := range events {
			switch event.Type {
			case linebot.EventTypeMessage:
				handleTextMessage(bot, event)
			case linebot.EventTypeFollow:
				handleGreet(bot, event)
				// Other event types...
			}
		}
		c.Status(http.StatusOK)
	}
}

func handleGreet(bot *linebot.Client, event *linebot.Event) {

	replyToken := event.ReplyToken
	// Get user profile
	profile, err := bot.GetProfile(event.Source.UserID).Do()
	if err != nil {
		log.Println("Error getting user profile:", err)
		sendReply(bot, replyToken, "Failed to get user profile. Please try again later.")
		return
	}

	// Extract user information from the profile
	replyText := "สวัสดีคุณ" + profile.DisplayName
	sendReply(bot, replyToken, replyText)
}

func handleTextMessage(bot *linebot.Client, event *linebot.Event) {
	replyToken := event.ReplyToken
	message := event.Message.(*linebot.TextMessage)
	text := message.Text
	switch text {
	case "reg":
		// Get user profile
		profile, err := bot.GetProfile(event.Source.UserID).Do()
		if err != nil {
			log.Println("Error getting user profile:", err)
			sendReply(bot, replyToken, "Failed to get user profile. Please try again later.")
			return
		}

		// Extract user information from the profile
		lineID := profile.UserID
		displayName := profile.DisplayName
		pictureURL := profile.PictureURL

		// Prepare SQL statement
		query := `
            INSERT INTO pokerth.member (
                member_id,
                member_line_Id,
                member_display_Name,
                member_credit,
                userLevelId,
                member_picture_url
            ) VALUES (
				?, ?, ?, ?, ?, ?, ?
			)
        `

		// Execute the SQL statement
		_, err = models.GetDB().Exec(query, "", lineID, displayName, 0, 4, pictureURL)
		if err != nil {
			log.Println("Error inserting new member:", err)
			sendReply(bot, replyToken, "Failed to register. Please try again later.")
			return
		}

		// Registration successful message
		sendReply(bot, replyToken, "Registration successful!")

	case "showId":
		// Implement function to retrieve user's Line ID
		sendReply(bot, replyToken, "Your Line ID is: "+event.Source.UserID)
	case "c", "C":
		// Implement function to retrieve user's credit from database and reply
		var userCredit float64
		err := models.GetDB().QueryRow("SELECT member_credit FROM member WHERE member_line_Id = ?", event.Source.UserID).Scan(&userCredit)
		if err != nil {
			log.Println("Error retrieving user credit from database:", err)
			sendReply(bot, replyToken, "Error retrieving user credit from database")
			return
		}

		sendReply(bot, replyToken, "Your credit is: "+strconv.FormatFloat(userCredit, 'f', 2, 64))
	default:
		// Handle other commands or messages
	}
}

func sendReply(bot *linebot.Client, replyToken, message string) {
	reply := linebot.NewTextMessage(message)
	_, err := bot.ReplyMessage(replyToken, reply).Do()
	if err != nil {
		log.Print("Error replying to user:", err)
	}
}
