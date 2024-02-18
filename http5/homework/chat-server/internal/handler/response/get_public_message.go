package response

import "time"

type GetPublicMessageResponse struct {
	FromUsername string    `json:"from_username"`
	Content      string    `json:"content"`
	SentAt       time.Time `json:"sent_at"`
	EditedAt     time.Time `json:"edited_at"`
}
