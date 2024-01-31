package service

import (
	"context"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/dto"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/model"
	sliceutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/slice"
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

func (ms *MessageService) SendPrivateMessage(ctx context.Context, createModel dto.CreatePrivateMessageDTO) (*model.PrivateMessage, error) {
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

func (ms *MessageService) SendPublicMessage(ctx context.Context, createModel dto.CreatePublicMessageDTO) (*model.PublicMessage, error) {
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
	// todo: maybe we should validate that user that requests this message is a sender or receiver?
	msg, err := ms.PrivateMessageRepo.GetPrivateMessage(ctx, id)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (ms *MessageService) GetAllPrivateMessages(ctx context.Context, userFrom *model.User) []*model.PrivateMessage {
	messages := ms.PrivateMessageRepo.GetAllPrivateMessages(ctx)

	// return only messages that were sent to current user
	return sliceutils.Filter(messages, func(msg *model.PrivateMessage) bool { return msg.To.ID == userFrom.ID })
}

func (ms *MessageService) GetAllPrivateMessagesFromUser(ctx context.Context, user *model.User, id int) []*model.PrivateMessage {
	messages := ms.PrivateMessageRepo.GetAllPrivateMessages(ctx)

	return sliceutils.Filter(messages, func(msg *model.PrivateMessage) bool { return msg.From.ID == id && msg.To.ID == user.ID })
}

func (ms *MessageService) GetPublicMessage(ctx context.Context, id int) (*model.PublicMessage, error) {
	msg, err := ms.PublicMessageRepo.GetPublicMessage(ctx, id)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (ms *MessageService) GetAllPublicMessages(ctx context.Context) []*model.PublicMessage {
	return ms.PublicMessageRepo.GetAllPublicMessages(ctx)
}

func (ms *MessageService) UpdatePrivateMessage(ctx context.Context, id int, newContent string) (*model.PrivateMessage, error) {
	msg, err := ms.PrivateMessageRepo.UpdatePrivateMessage(ctx, id, newContent)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (ms *MessageService) UpdatePublicMessage(ctx context.Context, id int, newContent string) (*model.PublicMessage, error) {
	msg, err := ms.PublicMessageRepo.UpdatePublicMessage(ctx, id, newContent)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (ms *MessageService) DeletePrivateMessage(ctx context.Context, id int) (*model.PrivateMessage, error) {
	msg, err := ms.PrivateMessageRepo.DeletePrivateMessage(ctx, id)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (ms *MessageService) DeletePublicMessage(ctx context.Context, id int) (*model.PublicMessage, error) {
	msg, err := ms.PublicMessageRepo.DeletePublicMessage(ctx, id)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
