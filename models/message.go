package models

type Message struct {
	UserID         string  `json:"user_id"`
	LineID         string  `json:"line_id"`
	LineName       string  `json:"line_name"`
	LineDisplayUrl string  `json:"line_display_url"`
	Credit         float64 `json:"credit"`
}
