package models

type MessageResponse struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type ReplyMessage struct {
	ReplyToken string            `json:"replyToken"`
	Messages   []MessageResponse `json:"messages"`
}
