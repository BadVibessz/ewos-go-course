package user

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entity"
)

//go:generate mockgen -destination=../../mocks/hasher.go -package=mocks github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/service/user Hasher

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

type Hasher interface {
	GenerateFromPassword(password []byte, cost int) ([]byte, error)
}

type Service struct {
	UserRepo UserRepo
	Hasher   Hasher
}

func New(userRepo UserRepo, hasher Hasher) *Service {
	return &Service{
		UserRepo: userRepo,
		Hasher:   hasher,
	}
}

func (us *Service) RegisterUser(ctx context.Context, user entity.User) (*entity.User, error) {
	// ensure that user with this email and username does not exist
	err := us.UserRepo.CheckUniqueConstraints(ctx, user.Email, user.Username)
	if err != nil {
		return nil, err
	}

	// user model sent with plain password
	hash, err := us.Hasher.GenerateFromPassword([]byte(user.HashedPassword), bcrypt.DefaultCost)
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

func initEmptyFieldsOfUser(usr1, usr2 *entity.User) {
	if usr1.Email == "" {
		usr1.Email = usr2.Email
	}

	if usr1.Username == "" {
		usr1.Username = usr2.Username
	}

	if usr1.HashedPassword == "" {
		usr1.HashedPassword = usr2.HashedPassword
	}
}

func usersEquals(usr1, usr2 *entity.User) bool {
	return usr1.Email == usr2.Email &&
		usr1.Username == usr2.Username &&
		usr1.HashedPassword == usr2.HashedPassword
}

func (us *Service) UpdateUser(ctx context.Context, id int, updateModel entity.User) (*entity.User, error) {
	err := us.UserRepo.CheckUniqueConstraints(ctx, updateModel.Email, updateModel.Username)
	if err != nil {
		return nil, err
	}

	usr, err := us.UserRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// if password changed => hash
	if updateModel.HashedPassword != "" {
		hash, err := us.Hasher.GenerateFromPassword([]byte(updateModel.HashedPassword), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		updateModel.HashedPassword = string(hash)
	}

	initEmptyFieldsOfUser(&updateModel, usr)

	// if update model updates nothing, thus no need for UserRepo.UpdateUser() call
	if usersEquals(&updateModel, usr) {
		return usr, nil
	}

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
