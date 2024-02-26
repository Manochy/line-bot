package models

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "root:Evild0ergu@tcp(localhost:3306)/pokerth")
	if err != nil {
		log.Fatal(err)
	}
}

func GetDB() *sql.DB {
	return db
}

// Member represents the structure of a member in the database
type Member struct {
	MemberID          string  `json:"member_id"`
	MemberLineID      string  `json:"member_line_Id"`
	MemberDisplayName string  `json:"member_display_Name"`
	MemberCredit      float64 `json:"member_credit"`
	UserLevelID       int     `json:"userLevelId"`
}
