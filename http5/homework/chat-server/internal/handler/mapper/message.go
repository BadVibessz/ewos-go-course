package mapper

import (
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entity"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/request"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/response"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/model"
)

func MapPublicMessageToResponse(msg *model.PublicMessage) response.PublicMessageResponse {
	return response.PublicMessageResponse{
		FromUsername: msg.From.Username,
		Content:      msg.Content,
		SentAt:       msg.SentAt,
		EditedAt:     msg.EditedAt,
	}
}

func MapPublicMessageRequestToEntity(req *request.SendPublicMessageRequest) entity.PublicMessage {
	return entity.PublicMessage{
		FromID:  req.FromID,
		Content: req.Content,
	}
}

func MapPrivateMessageToResponse(msg *model.PrivateMessage) response.PrivateMessageResponse {
	return response.PrivateMessageResponse{
		FromUsername: msg.From.Username,
		ToUsername:   msg.To.Username,
		Content:      msg.Content,
		SentAt:       msg.SentAt,
		EditedAt:     msg.EditedAt,
	}
}

func MapPrivateMessageRequestToEntity(req *request.SendPrivateMessageRequest) entity.PrivateMessage {
	return entity.PrivateMessage{
		ToID:    req.ToID,
		FromID:  req.FromID,
		Content: req.Content,
	}
}
