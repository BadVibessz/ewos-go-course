package user

import (
	"context"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/dto"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/model"

	usermapper "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/pkg/mapper/user"
)

type UserRepo interface {
	AddUser(ctx context.Context, user model.User) (*model.User, error)
	GetUserByID(ctx context.Context, id int) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetAllUsers(ctx context.Context, offset, limit int) []*model.User
	DeleteUser(ctx context.Context, id int) (*model.User, error)
	UpdateUser(ctx context.Context, id int, updateModel model.User) (*model.User, error)
	CheckUniqueConstraints(ctx context.Context, email, username string) error
}

type UserService struct {
	UserRepo UserRepo
}

func NewUserService(ur UserRepo) *UserService {
	return &UserService{UserRepo: ur}
}

func (us *UserService) RegisterUser(ctx context.Context, user dto.UserDTO) (*model.User, error) {
	// ensure that user with this email and username does not exist
	err := us.UserRepo.CheckUniqueConstraints(ctx, user.Email, user.Username)
	if err != nil {
		return nil, err
	}

	created, err := us.UserRepo.AddUser(ctx, usermapper.MapUserDtoToUser(&user))
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

func (us *UserService) GetAllUsers(ctx context.Context, offset, limit int) []*model.User {
	return us.UserRepo.GetAllUsers(ctx, offset, limit)
}

func (us *UserService) UpdateUser(ctx context.Context, id int, updateModel dto.UserDTO) (*model.User, error) {
	user := usermapper.MapUserDtoToUser(&updateModel)

	updated, err := us.UserRepo.UpdateUser(ctx, id, user)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (us *UserService) DeleteUser(ctx context.Context, id int) (*model.User, error) { // todo: authorize admin rights
	deleted, err := us.UserRepo.DeleteUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return deleted, nil
}
