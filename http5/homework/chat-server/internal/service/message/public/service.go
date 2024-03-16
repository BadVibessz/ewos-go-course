package public

import (
	"context"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entity"
)

//go:generate mockgen -destination=../../../mocks/public_message_repository.go -package=mocks github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/service/message/public PublicMessageRepo

type PublicMessageRepo interface {
	AddPublicMessage(ctx context.Context, msg entity.PublicMessage) (*entity.PublicMessage, error)
	GetAllPublicMessages(ctx context.Context, offset, limit int) []*entity.PublicMessage
	GetPublicMessage(ctx context.Context, id int) (*entity.PublicMessage, error)
}

type UserRepo interface {
	GetAllUsers(ctx context.Context, offset, limit int) []*entity.User
	GetUserByID(ctx context.Context, id int) (*entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
}

type Service struct {
	PublicMessageRepo PublicMessageRepo
	UserRepo          UserRepo
}

func New(publicMessageRepo PublicMessageRepo, userRepo UserRepo) *Service {
	return &Service{
		PublicMessageRepo: publicMessageRepo,
		UserRepo:          userRepo,
	}
}

func (s *Service) SendPublicMessage(ctx context.Context, msg entity.PublicMessage) (*entity.PublicMessage, error) {
	// check if user with provided username exists in database
	if _, err := s.UserRepo.GetUserByUsername(ctx, msg.FromUsername); err != nil {
		return nil, err
	}

	created, err := s.PublicMessageRepo.AddPublicMessage(ctx, msg)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *Service) GetPublicMessage(ctx context.Context, id int) (*entity.PublicMessage, error) {
	msg, err := s.PublicMessageRepo.GetPublicMessage(ctx, id)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (s *Service) GetAllPublicMessages(ctx context.Context, offset, limit int) []*entity.PublicMessage {
	return s.PublicMessageRepo.GetAllPublicMessages(ctx, offset, limit)
}
