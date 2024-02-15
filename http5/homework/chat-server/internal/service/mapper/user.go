package mapper

import (
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entity"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/model"
)

func MapUserDtoToUser(user *entity.User) model.User {
	return model.User{
		Email:          user.Email,
		Username:       user.Username,
		HashedPassword: user.HashedPassword,
	}
}
