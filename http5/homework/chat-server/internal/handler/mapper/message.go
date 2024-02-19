package mapper

import (
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entity"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/response"
)

func MapPublicMessageToResponse(msg *entity.PublicMessage) response.GetPublicMessageResponse {
	return response.GetPublicMessageResponse{
		FromUsername: msg.From.Username,
		Content:      msg.Content,
		SentAt:       msg.SentAt,
		EditedAt:     msg.EditedAt,
	}
}

func MapPrivateMessageToResponse(msg *entity.PrivateMessage) response.GetPrivateMessageResponse {
	return response.GetPrivateMessageResponse{
		FromUsername: msg.From.Username,
		ToUsername:   msg.To.Username,
		Content:      msg.Content,
		SentAt:       msg.SentAt,
		EditedAt:     msg.EditedAt,
	}
}
