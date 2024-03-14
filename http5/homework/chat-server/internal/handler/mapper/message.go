package mapper

import (
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entity"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/request"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/response"
)

func MapPublicMessageToResponse(msg *entity.PublicMessage) response.GetPublicMessageResponse {
	return response.GetPublicMessageResponse{
		FromUsername: msg.FromUsername,
		Content:      msg.Content,
		SentAt:       msg.SentAt,
		EditedAt:     msg.EditedAt,
	}
}

func MapPrivateMessageToResponse(msg *entity.PrivateMessage) response.GetPrivateMessageResponse {
	return response.GetPrivateMessageResponse{
		FromUsername: msg.FromUsername,
		ToUsername:   msg.ToUsername,
		Content:      msg.Content,
		SentAt:       msg.SentAt,
		EditedAt:     msg.EditedAt,
	}
}

func MapSendPrivateMessageRequestToEntity(req request.SendPrivateMessageRequest, fromUsername string) entity.PrivateMessage {
	return entity.PrivateMessage{
		FromUsername: fromUsername,
		ToUsername:   req.ToUsername,
		Content:      req.Content,
	}
}

func MapSendPublicMessageRequestToEntity(req request.SendPublicMessageRequest, fromUsername string) entity.PublicMessage {
	return entity.PublicMessage{
		FromUsername: fromUsername,
		Content:      req.Content,
	}
}
