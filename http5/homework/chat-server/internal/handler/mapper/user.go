package mapper

import (
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entity"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/request"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/response"
)

func MapUserToUserResponse(user *entity.User) response.GetUserResponse {
	return response.GetUserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func MapRegisterRequestToUserEntity(registerReq *request.RegisterRequest) entity.User {
	return entity.User{
		Email:          registerReq.Email,
		Username:       registerReq.Username,
		HashedPassword: registerReq.Password,
	}
}
