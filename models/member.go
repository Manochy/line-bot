package models

import (
	"strconv"

	"gorm.io/gorm"
)

// Member represents the structure of the member table
type Member struct {
	MemberID          string `gorm:"primaryKey"`
	MemberLineID      string
	MemberDisplayName string
	MemberCredit      float64
	UserLevelID       int
	MemberPictureURL  string `gorm:"type:VARCHAR(255)"`
}

func (Member) TableName() string {
	return "member"
}

// CreateMember inserts a new member record into the database
func CreateMember(db *gorm.DB, member *Member) error {
	return db.Create(member).Error
}

// UpdateMember updates an existing member record in the database
func UpdateMember(db *gorm.DB, member *Member) error {
	return db.Save(member).Error
}

// DeleteMember deletes an existing member record from the database
func DeleteMember(db *gorm.DB, memberID string) error {
	return db.Delete(&Member{}, "member_id = ?", memberID).Error
}

// SelectMember retrieves a member record from the database based on member ID
func SelectMemberByLineId(db *gorm.DB, memberID string) (*Member, error) {
	var member Member
	err := db.First(&member, "member_line_Id = ?", memberID).Error
	return &member, err
}

// GetNextMemberID returns the next available MemberID by incrementing the maximum MemberID in the database
func GetNextMemberID(db *gorm.DB) (string, error) {
	var maxMemberID string
	if err := db.Raw("SELECT MAX(member_id) FROM member").Scan(&maxMemberID).Error; err != nil {
		return "", err
	}

	// Convert the maximum MemberID to an integer
	maxID, err := strconv.Atoi(maxMemberID)
	if err != nil {
		return "", err
	}

	// Increment the maximum MemberID by 1
	nextMemberID := strconv.Itoa(maxID + 1)

	return nextMemberID, nil
}
