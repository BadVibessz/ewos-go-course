package private

import (
	"context"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/service/message"
	"math"
	"slices"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entity"

	sliceutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/slice"
)

//go:generate mockgen -destination=../../../mocks/private_message_repository.go -package=mocks github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/service/message/private PrivateMessageRepo

type PrivateMessageRepo interface {
	AddPrivateMessage(ctx context.Context, msg entity.PrivateMessage) (*entity.PrivateMessage, error)
	GetAllPrivateMessages(ctx context.Context, offset, limit int) []*entity.PrivateMessage
	GetPrivateMessage(ctx context.Context, id int) (*entity.PrivateMessage, error)
}

type UserRepo interface {
	GetAllUsers(ctx context.Context, offset, limit int) []*entity.User
	GetUserByID(ctx context.Context, id int) (*entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
}

type Service struct {
	PrivateMessageRepo PrivateMessageRepo
	UserRepo           UserRepo
}

func New(privateMessageRepo PrivateMessageRepo, userRepo UserRepo) *Service {
	return &Service{
		PrivateMessageRepo: privateMessageRepo,
		UserRepo:           userRepo,
	}
}

func (s *Service) checkSenderAndReceiver(ctx context.Context, senderUsername, receiverUsername string) error {
	if _, err := s.UserRepo.GetUserByUsername(ctx, senderUsername); err != nil {
		return message.ErrNoSuchReceiver
	}

	if _, err := s.UserRepo.GetUserByUsername(ctx, receiverUsername); err != nil {
		return message.ErrNoSuchReceiver
	}

	return nil
}

func (s *Service) SendPrivateMessage(ctx context.Context, msg entity.PrivateMessage) (*entity.PrivateMessage, error) {
	// check if users with provided usernames exists in database
	if err := s.checkSenderAndReceiver(ctx, msg.FromUsername, msg.ToUsername); err != nil {
		return nil, err
	}

	created, err := s.PrivateMessageRepo.AddPrivateMessage(ctx, msg)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *Service) GetPrivateMessage(ctx context.Context, id int) (*entity.PrivateMessage, error) {
	// todo: we should validate that user that requests this message is a sender or receiver
	msg, err := s.PrivateMessageRepo.GetPrivateMessage(ctx, id)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (s *Service) GetAllPrivateMessages(ctx context.Context, toUsername string, offset, limit int) []*entity.PrivateMessage {
	messages := s.PrivateMessageRepo.GetAllPrivateMessages(ctx, 0, math.MaxInt64)

	// return only messages that were sent to or from current user
	messages = sliceutils.Filter(messages, func(msg *entity.PrivateMessage) bool {
		return msg.ToUsername == toUsername || msg.FromUsername == toUsername
	})

	return sliceutils.Slice(messages, offset, limit)
}

func (s *Service) GetAllPrivateMessagesFromUser(ctx context.Context, toUsername, fromUsername string, offset, limit int) ([]*entity.PrivateMessage, error) {
	if err := s.checkSenderAndReceiver(ctx, fromUsername, toUsername); err != nil {
		return nil, err
	}

	messages := s.PrivateMessageRepo.GetAllPrivateMessages(ctx, 0, math.MaxInt64)
	messages = sliceutils.Filter(messages, func(msg *entity.PrivateMessage) bool {
		return msg.FromUsername == fromUsername && msg.ToUsername == toUsername
	})

	return sliceutils.Slice(messages, offset, limit), nil
}

func (s *Service) GetAllUsersThatSentMessage(ctx context.Context, toUsername string, offset, limit int) []*entity.User {
	messages := s.GetAllPrivateMessages(ctx, toUsername, offset, limit)
	usersFromIds := sliceutils.Unique(sliceutils.Map(messages, func(msg *entity.PrivateMessage) string { return msg.FromUsername }))

	allUsers := s.UserRepo.GetAllUsers(ctx, 0, math.MaxInt64)

	res := make([]*entity.User, 0, len(usersFromIds))
	for _, usr := range allUsers {
		if slices.Contains(usersFromIds, usr.Username) {
			res = append(res, usr)
		}
	}

	return res
}
