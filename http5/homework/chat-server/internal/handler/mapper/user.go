package mapper

import (
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entity"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/request"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/response"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/model"
)

func MapUserToUserResponse(user *model.User) response.UserResponse {
	return response.UserResponse{
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
