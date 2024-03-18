package auth

import (
	"context"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entity"
)

type UserRepo interface {
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
}

//go:generate mockgen -destination=../../mocks/hasher_auth.go -package=mocks github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/service/auth Hasher

type Hasher interface {
	CompareHashAndPassword(hashedPassword []byte, password []byte) error
}

type Service struct {
	UserRepo UserRepo
	Hasher   Hasher
}

func New(ur UserRepo, hasher Hasher) *Service {
	return &Service{
		UserRepo: ur,
		Hasher:   hasher,
	}
}

func (as *Service) Login(ctx context.Context, username, password string) (*entity.User, error) {
	user, err := as.UserRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	err = as.Hasher.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return nil, err
	}

	return user, nil
}
