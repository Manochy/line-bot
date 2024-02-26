package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Manochy/line-bot/models"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
	"gorm.io/gorm"
)

const failRegister = "ลงทะเบียนล้มเหลว กรุณาติดต่อแอดมินค่ะ"

func HandleCallback(bot *linebot.Client, db *gorm.DB) gin.HandlerFunc {
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
				handleTextMessage(bot, db, event)
			case linebot.EventTypeFollow:
				handleGreet(bot, db, event)
				// Other event types...
			}
		}
		c.Status(http.StatusOK)
	}
}

func handleGreet(bot *linebot.Client, db *gorm.DB, event *linebot.Event) {

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

func handleTextMessage(bot *linebot.Client, db *gorm.DB, event *linebot.Event) {
	//Init always used variable
	replyToken := event.ReplyToken
	message := event.Message.(*linebot.TextMessage)
	text := message.Text

	// Get user profile
	profile, err := bot.GetProfile(event.Source.UserID).Do()
	if err != nil {
		log.Println("Error getting user profile:", err)
		sendReply(bot, replyToken, "ดึงข้อมูลจากทางเซิฟเวอร์ของ Line ไม่ได้ \n โปรดรอสักครู่ค่อยทำรายการใหม่ค่ะ")
		return
	}

	// Extract user information from the profile
	lineID := profile.UserID

	switch text {
	case "reg":
		// Check if a member with the same line ID already exists
		existingMember, err := models.SelectMemberByLineId(db, lineID)
		if err != nil {
			log.Println("Error retrieving existing member:", err)
			sendReply(bot, replyToken, failRegister)
			return
		}

		// If a member with the same line ID already exists, do nothing
		if existingMember != nil {
			sendReply(bot, replyToken, "ไอดีคุณเคยลงทะเบียนไว้แล้วค่ะ \n ตรวจสอบเครดิตพิมพ์ c ค่ะ")
			return
		}

		// Generate the next MemberID
		nextMemberId, err := models.GetNextMemberID(db)
		if err != nil {
			log.Println("Error generating next MemberID:", err)
			sendReply(bot, replyToken, failRegister)
			return
		}

		member := models.Member{
			MemberID:          nextMemberId,
			MemberLineID:      lineID,
			MemberDisplayName: profile.DisplayName,
			MemberPictureURL:  profile.PictureURL,
			MemberCredit:      0.0,
			UserLevelID:       4,
		}
		err = models.CreateMember(db, &member)
		if err != nil {
			// Handle the error (e.g., return an error response)
			log.Println("Error inserting new member:", err)
			sendReply(bot, replyToken, failRegister)
			return
		}

		// Registration successful message
		sendReply(bot, replyToken, "ลงทะเบียนสำเร็จเรียบร้อยค่ะ")

	case "showId":
		// Implement function to retrieve user's Line ID
		sendReply(bot, replyToken, "ไอดีไลน์ของคุณคือ : "+event.Source.UserID)

	case "c", "C":
		// Implement function to retrieve user's credit from database and reply
		member, err := models.SelectMemberByLineId(db, lineID)
		if err != nil {
			// Handle the error (e.g., return an error response)
			log.Println("Error retrieving member:", err)
			sendReply(bot, replyToken, "Failed to retrieve member information. Please try again later.")
			return
		}

		// Get the credit from the retrieved member
		userCredit := member.MemberCredit

		sendReply(bot, replyToken, "คุณมีเครดิตคงเหลือ : "+strconv.FormatFloat(userCredit, 'f', 2, 64))
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
