package user

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entity"
)

type UserRepo interface {
	AddUser(ctx context.Context, user entity.User) (*entity.User, error)
	GetUserByID(ctx context.Context, id int) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
	GetAllUsers(ctx context.Context, offset, limit int) []*entity.User
	DeleteUser(ctx context.Context, id int) (*entity.User, error)
	UpdateUser(ctx context.Context, id int, updateModel entity.User) (*entity.User, error)
	CheckUniqueConstraints(ctx context.Context, email, username string) error
}

type Service struct {
	UserRepo UserRepo
}

func New(ur UserRepo) *Service {
	return &Service{UserRepo: ur}
}

func (us *Service) RegisterUser(ctx context.Context, user entity.User) (*entity.User, error) {
	// ensure that user with this email and username does not exist
	err := us.UserRepo.CheckUniqueConstraints(ctx, user.Email, user.Username)
	if err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.HashedPassword), bcrypt.DefaultCost) // user model sent with plain password
	if err != nil {
		return nil, err
	}

	user.HashedPassword = string(hash)

	created, err := us.UserRepo.AddUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (us *Service) GetUserByID(ctx context.Context, id int) (*entity.User, error) {
	user, err := us.UserRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *Service) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := us.UserRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *Service) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	user, err := us.UserRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *Service) GetAllUsers(ctx context.Context, offset, limit int) []*entity.User {
	return us.UserRepo.GetAllUsers(ctx, offset, limit)
}

func (us *Service) UpdateUser(ctx context.Context, id int, updateModel entity.User) (*entity.User, error) {
	updated, err := us.UserRepo.UpdateUser(ctx, id, updateModel)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (us *Service) DeleteUser(ctx context.Context, id int) (*entity.User, error) { // todo: authorize admin rights
	deleted, err := us.UserRepo.DeleteUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return deleted, nil
}
