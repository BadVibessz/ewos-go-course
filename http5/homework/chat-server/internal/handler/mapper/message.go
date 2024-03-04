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

func MapSendPrivateMessageRequestToEntity(req request.SendPrivateMessageRequest) entity.PrivateMessage {
	return entity.PrivateMessage{
		FromUsername: req.FromUsername,
		ToUsername:   req.ToUsername,
		Content:      req.Content,
	}
}

func MapSendPublicMessageRequestToEntity(req request.SendPublicMessageRequest) entity.PublicMessage {
	return entity.PublicMessage{
		FromUsername: req.FromUsername,
		Content:      req.Content,
	}
}
