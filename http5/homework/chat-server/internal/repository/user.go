package repository

import (
	"context"
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/model"

	inmemory "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/db/in-memory"
	sliceutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/slice"
)

const userTableName = "users"

var (
	ErrNoSuchUser     = errors.New("no such user")
	ErrEmailExists    = errors.New("user with this email already exists")
	ErrUsernameExists = errors.New("user with this username already exists")
)

type UserRepoInMemDB struct {
	DB inmemory.InMemoryDB
}

func NewInMemUserRepo(db inmemory.InMemoryDB) *UserRepoInMemDB {
	repo := UserRepoInMemDB{DB: db}

	repo.DB.CreateTable(userTableName)

	return &repo
}

func (ur *UserRepoInMemDB) GetAllUsers(_ context.Context, offset, limit int) []*model.User {
	rows, err := ur.DB.GetAllRows(userTableName, offset, limit)
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

func (ur *UserRepoInMemDB) AddUser(_ context.Context, user model.User) (*model.User, error) {
	idOffset, err := ur.DB.GetRowsCount(userTableName)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	user.ID = idOffset + 1
	user.CreatedAt = now
	user.UpdatedAt = now

	err = ur.DB.AddRow(userTableName, strconv.Itoa(user.ID), user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepoInMemDB) GetUserByID(_ context.Context, id int) (*model.User, error) {
	row, err := ur.DB.GetRow(userTableName, strconv.Itoa(id))
	if err != nil {
		return nil, ErrNoSuchUser
	}

	user, ok := row.(model.User)
	if !ok {
		return nil, ErrNoSuchUser
	}

	return &user, nil
}

func (ur *UserRepoInMemDB) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	users := ur.GetAllUsers(ctx, 0, math.MaxInt64)
	if len(users) == 0 {
		return nil, ErrNoSuchUser
	}

	filtered := sliceutils.Filter(users, func(u *model.User) bool { return u.Email == email })
	if len(filtered) == 0 {
		return nil, ErrNoSuchUser
	}

	user := filtered[0]

	return user, nil
}

func (ur *UserRepoInMemDB) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	users := ur.GetAllUsers(ctx, 0, math.MaxInt64)
	if len(users) == 0 {
		return nil, ErrNoSuchUser
	}

	filtered := sliceutils.Filter(users, func(u *model.User) bool { return u.Username == username })

	if len(filtered) == 0 {
		return nil, ErrNoSuchUser
	}

	user := filtered[0]

	return user, nil
}

func (ur *UserRepoInMemDB) DeleteUser(ctx context.Context, id int) (*model.User, error) {
	user, err := ur.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = ur.DB.DropRow(userTableName, strconv.Itoa(id))
	if err != nil {
		return nil, ErrNoSuchUser
	}

	return user, nil
}

func (ur *UserRepoInMemDB) UpdateUser(ctx context.Context, id int, updated model.User) (*model.User, error) {
	user, err := ur.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	updated.ID = id
	updated.CreatedAt = user.CreatedAt
	updated.UpdatedAt = time.Now()

	err = ur.DB.AlterRow(userTableName, strconv.Itoa(id), updated)
	if err != nil {
		return nil, ErrNoSuchUser
	}

	return user, nil
}

func (ur *UserRepoInMemDB) CheckUniqueConstraints(ctx context.Context, email, username string) error {
	got, err := ur.GetUserByEmail(ctx, email)
	if got != nil || err == nil {
		return ErrEmailExists
	}

	got, err = ur.GetUserByUsername(ctx, username)
	if got != nil || err == nil {
		return ErrUsernameExists
	}

	return nil
}
