package auth

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entity"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/request"
)

type UserRepo interface {
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
}

type Service struct {
	UserRepo UserRepo
}

func New(ur UserRepo) *Service {
	return &Service{
		UserRepo: ur,
	}
}

func (as *Service) Login(ctx context.Context, loginReq request.LoginRequest) (*entity.User, error) {
	user, err := as.UserRepo.GetUserByUsername(ctx, loginReq.Username)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(loginReq.Password))
	if err != nil {
		return nil, err
	}

	return user, nil
}
