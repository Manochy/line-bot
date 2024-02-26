// main.go

package main

import (
	"log"
	"os"

	"github.com/Manochy/line-bot/handlers"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var bot *linebot.Client
var db *gorm.DB

func initDB() {
	// Initialize database connection
	var err error
	dsn := "root:Evild0ergu@tcp(localhost:3306)/pokerth?charset=utf8mb4&parseTime=True&loc=Local" // Replace "dbname" with your database name
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
}

func main() {
	var err error
	bot, err = linebot.New(
		"43d644a6edbf6bbcf7f30b29672ff6be",
		"3sPiJyEc0px9rXJiv80h5HAvoqTF+RF5pGqfIcsh/5vEZ8+x+b2E064u8Yt4f5o2QwFhRit8G9Rn034AKMG4Z6Bxiur0qX7w9mCuasUlYerYW9H1D4b8sFLjQIxyN9cPW5VzbdxS/3FuYJKONzM+ewdB04t89/1O/w1cDnyilFU=",
	)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize db connection
	initDB()
	// Defer closing of db
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatal(err)
		}
		sqlDB.Close()
	}()

	r := gin.Default()
	r.POST("/callback", handlers.HandleCallback(bot, db))

	port := os.Getenv("PORT")
	if port == "" {
		port = "12024"
	}
	log.Printf("Starting server on port %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
