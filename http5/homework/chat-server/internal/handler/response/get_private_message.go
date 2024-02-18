package response

import "time"

type GetPrivateMessageResponse struct {
	FromUsername string    `json:"from_username"`
	ToUsername   string    `json:"to_username"`
	Content      string    `json:"content"`
	SentAt       time.Time `json:"sent_at"`
	EditedAt     time.Time `json:"edited_at"`
}
