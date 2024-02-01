package service

import (
	"context"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/dto"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type AuthBasicService struct {
	UserRepo UserRepo
}

func NewBasicAuthService(ur UserRepo) *AuthBasicService {
	return &AuthBasicService{
		UserRepo: ur,
	}
}

func (as *AuthBasicService) Login(ctx context.Context, cred dto.LoginUserDTO) (*model.User, error) {
	user, err := as.UserRepo.GetUserByUsername(ctx, cred.Username)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(cred.Password))
	if err != nil {
		return nil, err
	}

	return user, nil
}
