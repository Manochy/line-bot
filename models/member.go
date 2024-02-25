package models

// Member represents the structure of a member in the database
type Member struct {
	MemberID          string  `json:"member_id"`
	MemberLineID      string  `json:"member_line_Id"`
	MemberDisplayName string  `json:"member_display_Name"`
	MemberCredit      float64 `json:"member_credit"`
	UserLevelID       int     `json:"userLevelId"`
}
