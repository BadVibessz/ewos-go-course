package repositories

import (
	"context"
	"errors"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/models"
	inmemory "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/db/in-memory"
	"strconv"
	"time"
)

const userTableName = "users"

var (
	NoSuchUserErr = errors.New("no such user")
)

type UserRepo interface {
	AddUser(ctx context.Context, user models.CreateUserModel) (*models.User, error)
	GetUserById(ctx context.Context, id int) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetAllUsers(ctx context.Context) []*models.User
	DeleteUser(ctx context.Context, id int) (*models.User, error)
	UpdateUser(ctx context.Context, id int, model models.UpdateUserModel) (*models.User, error)
}

func NewInMemRepo(db *inmemory.InMemDB) UserRepo {
	repo := UserRepoInMemDB{DB: db}

	repo.DB.CreateTable(userTableName)

	return &repo
}

type UserRepoInMemDB struct {
	DB *inmemory.InMemDB // todo: is it good that ChatDB declared in server pkg?
}

func (ur *UserRepoInMemDB) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (ur *UserRepoInMemDB) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (ur *UserRepoInMemDB) GetAllUsers(ctx context.Context) []*models.User {
	rows, err := ur.DB.GetAllRows(userTableName)
	if err != nil {
		return nil
	}

	res := make([]*models.User, len(rows))
	for _, row := range rows {
		user, ok := row.(models.User)
		if ok {
			res = append(res, &user)
		}
	}

	return res
}

func (ur *UserRepoInMemDB) AddUser(ctx context.Context, user models.CreateUserModel) (*models.User, error) { // todo: how to use context?
	idOffset, err := ur.DB.GetRowsCount(userTableName)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	toCreate := models.User{
		ID:             idOffset + 1, // todo: user uuid.new()?
		Email:          user.Email,
		Username:       user.Username,
		HashedPassword: user.HashedPassword,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	err = ur.DB.AddRow(userTableName, strconv.Itoa(toCreate.ID), toCreate)
	if err != nil {
		return nil, err // todo: maybe return custom err of this layer? (NoSuchUserErr?)
	}

	return &toCreate, nil
}

func (ur *UserRepoInMemDB) GetUserById(ctx context.Context, id int) (*models.User, error) {
	row, err := ur.DB.GetRow(userTableName, strconv.Itoa(id))
	if err != nil {
		return nil, err
	}

	user, ok := row.(models.User)
	if !ok {
		return nil, NoSuchUserErr
	}

	return &user, nil
}

func (ur *UserRepoInMemDB) DeleteUser(ctx context.Context, id int) (*models.User, error) {
	user, err := ur.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	err = ur.DB.DropRow(userTableName, strconv.Itoa(id))
	if err != nil {
		return nil, NoSuchUserErr
	}

	return user, nil
}

func (ur *UserRepoInMemDB) UpdateUser(ctx context.Context, id int, updated models.UpdateUserModel) (*models.User, error) {
	user, err := ur.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	user.Username = updated.NewUsername
	user.Email = updated.NewEmail
	user.HashedPassword = updated.NewHashedPassword
	user.UpdatedAt = time.Now()

	err = ur.DB.AlterRow(userTableName, strconv.Itoa(id), user)
	if err != nil {
		return nil, NoSuchUserErr
	}

	return user, nil
}
