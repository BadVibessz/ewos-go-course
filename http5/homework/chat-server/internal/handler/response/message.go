package response

import (
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/model"
	"time"
)

type PublicMessageResponse struct {
	FromUsername string    `json:"from_username"`
	Content      string    `json:"content"`
	SentAt       time.Time `json:"sent_at"`
	EditedAt     time.Time `json:"edited_at"`
}

func PublicMsgRespFromMessage(msg model.PublicMessage) PublicMessageResponse {
	return PublicMessageResponse{
		FromUsername: msg.From.Username,
		Content:      msg.Content,
		SentAt:       msg.SentAt,
		EditedAt:     msg.EditedAt,
	}
}

type PrivateMessageResponse struct {
	FromUsername string    `json:"from_username"`
	ToUsername   string    `json:"to_username"`
	Content      string    `json:"content"`
	SentAt       time.Time `json:"sent_at"`
	EditedAt     time.Time `json:"edited_at"`
}

func PrivateMsgRespFromMessage(msg model.PrivateMessage) PrivateMessageResponse {
	return PrivateMessageResponse{
		FromUsername: msg.From.Username,
		ToUsername:   msg.To.Username,
		Content:      msg.Content,
		SentAt:       msg.SentAt,
		EditedAt:     msg.EditedAt,
	}
}
