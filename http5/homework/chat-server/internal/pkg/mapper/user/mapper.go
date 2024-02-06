package user

import (
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/dto"
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

func MapUserDtoToUser(user *dto.UserDTO) model.User {
	return model.User{
		Email:          user.Email,
		Username:       user.Username,
		HashedPassword: user.HashedPassword,
	}
}
