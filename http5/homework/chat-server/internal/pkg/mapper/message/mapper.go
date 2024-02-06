package message

import (
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/response"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/model"
)

func MapPublicMessageToPublicMsgResp(msg *model.PublicMessage) response.PublicMessageResponse {
	return response.PublicMessageResponse{
		FromUsername: msg.From.Username,
		Content:      msg.Content,
		SentAt:       msg.SentAt,
		EditedAt:     msg.EditedAt,
	}
}

func MapPrivateMessageToPrivateMsgResp(msg *model.PrivateMessage) response.PrivateMessageResponse {
	return response.PrivateMessageResponse{
		FromUsername: msg.From.Username,
		ToUsername:   msg.To.Username,
		Content:      msg.Content,
		SentAt:       msg.SentAt,
		EditedAt:     msg.EditedAt,
	}
}
