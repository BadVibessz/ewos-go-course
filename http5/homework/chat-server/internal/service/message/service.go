package message

import (
	"context"
	"errors"
	"math"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entity"
	sliceutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/slice"
)

type PrivateMessageRepo interface {
	AddPrivateMessage(ctx context.Context, msg entity.PrivateMessage) (*entity.PrivateMessage, error)
	GetAllPrivateMessages(ctx context.Context, offset, limit int) []*entity.PrivateMessage
	GetPrivateMessage(ctx context.Context, id int) (*entity.PrivateMessage, error)
}

type PublicMessageRepo interface {
	AddPublicMessage(ctx context.Context, msg entity.PublicMessage) (*entity.PublicMessage, error)
	GetAllPublicMessages(ctx context.Context, offset, limit int) []*entity.PublicMessage
	GetPublicMessage(ctx context.Context, id int) (*entity.PublicMessage, error)
}

type UserRepo interface {
	GetUserByID(ctx context.Context, id int) (*entity.User, error)
}

var (
	ErrNoSuchReceiver = errors.New("no such receiver")
	ErrNoSuchSender   = errors.New("no such sender")
)

type Service struct {
	PrivateMessageRepo PrivateMessageRepo
	PublicMessageRepo  PublicMessageRepo
	UserRepo           UserRepo
}

func New(pr PrivateMessageRepo, pb PublicMessageRepo, ur UserRepo) *Service {
	return &Service{
		PrivateMessageRepo: pr,
		PublicMessageRepo:  pb,
		UserRepo:           ur,
	}
}

func (ms *Service) SendPrivateMessage(ctx context.Context, fromID, toID int, content string) (*entity.PrivateMessage, error) {
	userFrom, err := ms.UserRepo.GetUserByID(ctx, fromID)
	if err != nil {
		return nil, ErrNoSuchSender
	}

	userTo, err := ms.UserRepo.GetUserByID(ctx, toID)
	if err != nil {
		return nil, ErrNoSuchReceiver
	}

	msg := entity.PrivateMessage{
		From:    userFrom,
		To:      userTo,
		Content: content,
	}

	created, err := ms.PrivateMessageRepo.AddPrivateMessage(ctx, msg)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (ms *Service) SendPublicMessage(ctx context.Context, fromID int, content string) (*entity.PublicMessage, error) {
	userFrom, err := ms.UserRepo.GetUserByID(ctx, fromID)
	if err != nil {
		return nil, err
	}

	msg := entity.PublicMessage{
		From:    userFrom,
		Content: content,
	}

	created, err := ms.PublicMessageRepo.AddPublicMessage(ctx, msg)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (ms *Service) GetPrivateMessage(ctx context.Context, id int) (*entity.PrivateMessage, error) {
	// todo: we should validate that user that requests this message is a sender or receiver
	msg, err := ms.PrivateMessageRepo.GetPrivateMessage(ctx, id)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (ms *Service) GetAllPrivateMessages(ctx context.Context, userToID int, offset, limit int) []*entity.PrivateMessage {
	messages := ms.PrivateMessageRepo.GetAllPrivateMessages(ctx, 0, math.MaxInt64)

	// return only messages that were sent to current user
	messages = sliceutils.Filter(messages, func(msg *entity.PrivateMessage) bool { return msg.To.ID == userToID })

	return sliceutils.Slice(messages, offset, limit)
}

func (ms *Service) GetAllPrivateMessagesFromUser(ctx context.Context, toID, fromID int, offset, limit int) ([]*entity.PrivateMessage, error) {
	_, err := ms.UserRepo.GetUserByID(ctx, fromID)
	if err != nil {
		return nil, err
	}

	messages := ms.PrivateMessageRepo.GetAllPrivateMessages(ctx, 0, math.MaxInt64)
	messages = sliceutils.Filter(messages, func(msg *entity.PrivateMessage) bool { return msg.From.ID == fromID && msg.To.ID == toID })

	return sliceutils.Slice(messages, offset, limit), nil
}

func (ms *Service) GetPublicMessage(ctx context.Context, id int) (*entity.PublicMessage, error) {
	msg, err := ms.PublicMessageRepo.GetPublicMessage(ctx, id)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (ms *Service) GetAllPublicMessages(ctx context.Context, offset, limit int) []*entity.PublicMessage {
	return ms.PublicMessageRepo.GetAllPublicMessages(ctx, offset, limit)
}

func (ms *Service) GetAllUsersThatSentMessage(ctx context.Context, toID int, offset, limit int) []*entity.User {
	messages := ms.GetAllPrivateMessages(ctx, toID, offset, limit)
	usersDuplicates := sliceutils.Map(messages, func(msg *entity.PrivateMessage) *entity.User { return msg.From })

	return sliceutils.Unique(usersDuplicates)
}
