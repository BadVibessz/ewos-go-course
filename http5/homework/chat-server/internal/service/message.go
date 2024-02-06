package service

import (
	"context"
	"math"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/dto"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/model"

	sliceutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/slice"
)

type PrivateMessageRepo interface {
	AddPrivateMessage(ctx context.Context, msg model.PrivateMessage) (*model.PrivateMessage, error)
	GetAllPrivateMessages(ctx context.Context, offset, limit int) []*model.PrivateMessage
	GetPrivateMessage(ctx context.Context, id int) (*model.PrivateMessage, error)
	UpdatePrivateMessage(ctx context.Context, id int, newContent string) (*model.PrivateMessage, error)
	DeletePrivateMessage(ctx context.Context, id int) (*model.PrivateMessage, error)
}

type PublicMessageRepo interface {
	AddPublicMessage(ctx context.Context, msg model.PublicMessage) (*model.PublicMessage, error)
	GetAllPublicMessages(ctx context.Context, offset, limit int) []*model.PublicMessage
	GetPublicMessage(ctx context.Context, id int) (*model.PublicMessage, error)
	UpdatePublicMessage(ctx context.Context, id int, newContent string) (*model.PublicMessage, error)
	DeletePublicMessage(ctx context.Context, id int) (*model.PublicMessage, error)
}

type UserRepoMsgService interface {
	AddUser(ctx context.Context, user model.User) (*model.User, error)
	GetUserByID(ctx context.Context, id int) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetAllUsers(ctx context.Context, offset, limit int) []*model.User
	DeleteUser(ctx context.Context, id int) (*model.User, error)
	UpdateUser(ctx context.Context, id int, updateModel model.User) (*model.User, error)
	CheckUniqueConstraints(ctx context.Context, email, username string) error
}

type MessageService struct {
	PrivateMessageRepo PrivateMessageRepo
	PublicMessageRepo  PublicMessageRepo
	UserRepo           UserRepoMsgService
}

func NewMessageService(pr PrivateMessageRepo, pb PublicMessageRepo, ur UserRepoMsgService) *MessageService {
	return &MessageService{
		PrivateMessageRepo: pr,
		PublicMessageRepo:  pb,
		UserRepo:           ur,
	}
}

func (ms *MessageService) SendPrivateMessage(ctx context.Context, createModel dto.PrivateMessageDTO) (*model.PrivateMessage, error) {
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

func (ms *MessageService) SendPublicMessage(ctx context.Context, createModel dto.PublicMessageDTO) (*model.PublicMessage, error) {
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

func (ms *MessageService) GetAllPrivateMessages(ctx context.Context, userToID int, offset, limit int) []*model.PrivateMessage {
	messages := ms.PrivateMessageRepo.GetAllPrivateMessages(ctx, 0, math.MaxInt64)

	// return only messages that were sent to current user
	messages = sliceutils.Filter(messages, func(msg *model.PrivateMessage) bool { return msg.To.ID == userToID })

	return sliceutils.Slice(messages, offset, limit)
}

func (ms *MessageService) GetAllPrivateMessagesFromUser(ctx context.Context, toID, fromID int, offset, limit int) []*model.PrivateMessage {
	messages := ms.PrivateMessageRepo.GetAllPrivateMessages(ctx, 0, math.MaxInt64)
	messages = sliceutils.Filter(messages, func(msg *model.PrivateMessage) bool { return msg.From.ID == fromID && msg.To.ID == toID })

	return sliceutils.Slice(messages, offset, limit)
}

func (ms *MessageService) GetPublicMessage(ctx context.Context, id int) (*model.PublicMessage, error) {
	msg, err := ms.PublicMessageRepo.GetPublicMessage(ctx, id)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (ms *MessageService) GetAllPublicMessages(ctx context.Context, offset, limit int) []*model.PublicMessage {
	return ms.PublicMessageRepo.GetAllPublicMessages(ctx, offset, limit)
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
