package response

import (
	"time"
)

type PublicMessageResponse struct {
	FromUsername string    `json:"from_username"`
	Content      string    `json:"content"`
	SentAt       time.Time `json:"sent_at"`
	EditedAt     time.Time `json:"edited_at"`
}

type PrivateMessageResponse struct {
	FromUsername string    `json:"from_username"`
	ToUsername   string    `json:"to_username"`
	Content      string    `json:"content"`
	SentAt       time.Time `json:"sent_at"`
	EditedAt     time.Time `json:"edited_at"`
}
