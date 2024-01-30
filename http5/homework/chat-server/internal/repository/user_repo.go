package repository

import (
	"context"
	"errors"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/model"
	inmemory "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/db/in-memory"
	sliceutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/slice"
	"strconv"
	"time"
)

const userTableName = "users"

var (
	NoSuchUserErr = errors.New("no such user")
)

func NewInMemUserRepo(db *inmemory.InMemDB) *UserRepoInMemDB {
	repo := UserRepoInMemDB{DB: db}

	repo.DB.CreateTable(userTableName)

	return &repo
}

type UserRepoInMemDB struct {
	DB *inmemory.InMemDB
}

func (ur *UserRepoInMemDB) GetAllUsers(ctx context.Context) []*model.User {
	rows, err := ur.DB.GetAllRows(userTableName)
	if err != nil {
		return nil
	}

	res := make([]*model.User, 0, len(rows))
	for _, row := range rows {
		user, ok := row.(model.User)
		if ok {
			res = append(res, &user)
		}
	}

	return res
}

func (ur *UserRepoInMemDB) AddUser(ctx context.Context, user model.User) (*model.User, error) { // todo: how to use context?
	idOffset, err := ur.DB.GetRowsCount(userTableName)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	toCreate := model.User{
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

func (ur *UserRepoInMemDB) GetUserByID(ctx context.Context, id int) (*model.User, error) {
	row, err := ur.DB.GetRow(userTableName, strconv.Itoa(id))
	if err != nil {
		return nil, err
	}

	user, ok := row.(model.User)
	if !ok {
		return nil, NoSuchUserErr
	}

	return &user, nil
}

func (ur *UserRepoInMemDB) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	users := ur.GetAllUsers(ctx)
	if len(users) == 0 {
		return nil, NoSuchUserErr
	}

	user := sliceutils.Filter(users, func(u *model.User) bool { return u.Email == email })[0]

	if user == nil {
		return nil, NoSuchUserErr
	}

	return user, nil
}

func (ur *UserRepoInMemDB) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	users := ur.GetAllUsers(ctx)
	if len(users) == 0 {
		return nil, NoSuchUserErr
	}

	user := sliceutils.Filter(users, func(u *model.User) bool { return u.Username == username })[0]

	if user == nil {
		return nil, NoSuchUserErr
	}

	return user, nil
}

func (ur *UserRepoInMemDB) DeleteUser(ctx context.Context, id int) (*model.User, error) {
	user, err := ur.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = ur.DB.DropRow(userTableName, strconv.Itoa(id))
	if err != nil {
		return nil, NoSuchUserErr
	}

	return user, nil
}

func (ur *UserRepoInMemDB) UpdateUser(ctx context.Context, id int, updated model.User) (*model.User, error) {
	user, err := ur.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.Username = updated.Username
	user.Email = updated.Email
	user.HashedPassword = updated.HashedPassword
	user.UpdatedAt = time.Now()

	err = ur.DB.AlterRow(userTableName, strconv.Itoa(id), user)
	if err != nil {
		return nil, NoSuchUserErr
	}

	return user, nil
}
