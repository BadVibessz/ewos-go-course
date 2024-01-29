package service

import (
	"context"
	"errors"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/dto"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/model"
)

type UserRepo interface {
	AddUser(ctx context.Context, user model.User) (*model.User, error)
	GetUserByID(ctx context.Context, id int) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetAllUsers(ctx context.Context) []*model.User
	DeleteUser(ctx context.Context, id int) (*model.User, error)
	UpdateUser(ctx context.Context, id int, model model.User) (*model.User, error)
}

var (
	UserWithThisEmailExistsErr    = errors.New("user with this email already exists")
	UserWithThisUsernameExistsErr = errors.New("user with this username already exists")
)

type UserService struct {
	UserRepo UserRepo
}

func NewUserService(ur UserRepo) *UserService {
	return &UserService{UserRepo: ur}
}

// TODO: BETTER ERRORS HANDLING (wrap errors)

func (us *UserService) RegisterUser(ctx context.Context, user dto.CreateUserDTO) (*model.User, error) {
	// ensure that user with this email and username does not exist
	got, _ := us.UserRepo.GetUserByEmail(ctx, user.Email)
	if got != nil {
		return nil, UserWithThisEmailExistsErr
	}

	got, _ = us.UserRepo.GetUserByUsername(ctx, user.Username)
	if got != nil {
		return nil, UserWithThisUsernameExistsErr
	}

	created, err := us.UserRepo.AddUser(ctx,
		model.User{
			Email:          user.Email,
			Username:       user.Username,
			HashedPassword: user.HashedPassword,
		})

	if err != nil {
		return nil, err
	}

	return created, nil
}

func (us *UserService) GetUserByID(ctx context.Context, id int) (*model.User, error) {
	user, err := us.UserRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := us.UserRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	user, err := us.UserRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) GetAllUsers(ctx context.Context) []*model.User {
	return us.UserRepo.GetAllUsers(ctx)
}

func (us *UserService) UpdateUser(ctx context.Context, id int, updateModel dto.UpdateUserDTO) (*model.User, error) {
	user := model.User{
		Email:          updateModel.NewEmail,
		Username:       updateModel.NewUsername,
		HashedPassword: updateModel.NewHashedPassword,
	}

	updated, err := us.UserRepo.UpdateUser(ctx, id, user)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (us *UserService) DeleteUser(ctx context.Context, id int) (*model.User, error) {
	deleted, err := us.UserRepo.DeleteUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return deleted, nil
}
