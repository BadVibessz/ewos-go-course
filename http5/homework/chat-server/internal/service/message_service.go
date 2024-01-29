package service

import (
	"context"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/dto"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/model"
)

type PrivateMessageRepo interface {
	AddPrivateMessage(ctx context.Context, msg model.PrivateMessage) (*model.PrivateMessage, error)
	GetAllPrivateMessages(ctx context.Context) []*model.PrivateMessage
	GetPrivateMessage(ctx context.Context, id int) (*model.PrivateMessage, error)
	UpdatePrivateMessage(ctx context.Context, id int, newContent string) (*model.PrivateMessage, error)
	DeletePrivateMessage(ctx context.Context, id int) (*model.PrivateMessage, error)
}

type PublicMessageRepo interface {
	AddPublicMessage(ctx context.Context, msg model.PublicMessage) (*model.PublicMessage, error)
	GetAllPublicMessages(ctx context.Context) []*model.PublicMessage
	GetPublicMessage(ctx context.Context, id int) (*model.PublicMessage, error)
	UpdatePublicMessage(ctx context.Context, id int, newContent string) (*model.PublicMessage, error)
	DeletePublicMessage(ctx context.Context, id int) (*model.PublicMessage, error)
}

type MessageService struct {
	PrivateMessageRepo PrivateMessageRepo
	PublicMessageRepo  PublicMessageRepo
	UserRepo           UserRepo
}

func NewMessageService(pr PrivateMessageRepo, pb PublicMessageRepo, ur UserRepo) *MessageService {
	return &MessageService{
		PrivateMessageRepo: pr,
		PublicMessageRepo:  pb,
		UserRepo:           ur,
	}
}

func (ms *MessageService) AddPrivateMessage(ctx context.Context, createModel dto.CreatePrivateMessageDTO) (*model.PrivateMessage, error) {
	userFrom, err := ms.UserRepo.GetUserByID(ctx, createModel.FromID)
	if err != nil {
		return nil, err
	}

	userTo, err := ms.UserRepo.GetUserByID(ctx, createModel.ToID)
	if err != nil {
		return nil, err
	}

	msg := model.PrivateMessage{
		From:    userFrom,
		To:      userTo,
		Content: createModel.Content,
	}

	created, err := ms.PrivateMessageRepo.AddPrivateMessage(ctx, msg)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (ms *MessageService) AddPublicMessage(ctx context.Context, createModel dto.CreatePublicMessageDTO) (*model.PublicMessage, error) {
	userFrom, err := ms.UserRepo.GetUserByID(ctx, createModel.FromID)
	if err != nil {
		return nil, err
	}

	msg := model.PublicMessage{
		From:    userFrom,
		Content: createModel.Content,
	}

	created, err := ms.PublicMessageRepo.AddPublicMessage(ctx, msg)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (ms *MessageService) GetPrivateMessage(ctx context.Context, id int) (*model.PrivateMessage, error) {
	msg, err := ms.PrivateMessageRepo.GetPrivateMessage(ctx, id)
	if err != nil {
		// todo:
	}

}
