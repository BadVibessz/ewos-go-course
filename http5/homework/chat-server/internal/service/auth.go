package service

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/model"

	userservice "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/service/user"
)

type AuthBasicService struct {
	UserRepo userservice.UserRepo
}

func NewBasicAuthService(ur userservice.UserRepo) *AuthBasicService {
	return &AuthBasicService{
		UserRepo: ur,
	}
}

func (as *AuthBasicService) Login(ctx context.Context, username, password string) (*model.User, error) {
	user, err := as.UserRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return nil, err
	}

	return user, nil
}
