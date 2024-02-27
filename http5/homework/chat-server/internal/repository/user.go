package repository

import (
	"context"
	"errors"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entity"

	inmemory "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/db/in-memory"
	sliceutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/slice"
)

type UserRepoInMemDB struct {
	DB    inmemory.InMemoryDB
	mutex sync.RWMutex
}

func NewInMemUserRepo(db inmemory.InMemoryDB) *UserRepoInMemDB {
	repo := UserRepoInMemDB{
		DB:    db,
		mutex: sync.RWMutex{},
	}

	_, err := repo.DB.GetTable(UserTableName)
	if err != nil && errors.Is(err, inmemory.ErrNotExistedTable) {
		repo.DB.CreateTable(UserTableName)
	}

	return &repo
}

func (ur *UserRepoInMemDB) GetAllUsers(_ context.Context, offset, limit int) []*entity.User {
	rows, err := ur.DB.GetAllRows(UserTableName, offset, limit)
	if err != nil {
		return nil
	}

	res := make([]*entity.User, 0, len(rows))

	for _, row := range rows {
		user, ok := row.(entity.User)
		if ok {
			res = append(res, &user)
		}
	}

	return res
}

func (ur *UserRepoInMemDB) AddUser(_ context.Context, user entity.User) (*entity.User, error) {
	ur.mutex.Lock()
	defer ur.mutex.Unlock()

	idOffset, err := ur.DB.GetTableCounter(UserTableName)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	user.ID = idOffset + 1
	user.CreatedAt = now
	user.UpdatedAt = now

	if err = ur.DB.AddRow(UserTableName, strconv.Itoa(user.ID), user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepoInMemDB) GetUserByID(_ context.Context, id int) (*entity.User, error) {
	row, err := ur.DB.GetRow(UserTableName, strconv.Itoa(id))
	if err != nil {
		return nil, ErrNoSuchUser
	}

	user, ok := row.(entity.User)
	if !ok {
		return nil, ErrNoSuchUser
	}

	return &user, nil
}

func (ur *UserRepoInMemDB) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	users := ur.GetAllUsers(ctx, 0, math.MaxInt64)
	if len(users) == 0 {
		return nil, ErrNoSuchUser
	}

	filtered := sliceutils.Filter(users, func(u *entity.User) bool { return u.Email == email })
	if len(filtered) == 0 {
		return nil, ErrNoSuchUser
	}

	user := filtered[0]

	return user, nil
}

func (ur *UserRepoInMemDB) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	users := ur.GetAllUsers(ctx, 0, math.MaxInt64)
	if len(users) == 0 {
		return nil, ErrNoSuchUser
	}

	filtered := sliceutils.Filter(users, func(u *entity.User) bool { return u.Username == username })

	if len(filtered) == 0 {
		return nil, ErrNoSuchUser
	}

	user := filtered[0]

	return user, nil
}

func (ur *UserRepoInMemDB) DeleteUser(ctx context.Context, id int) (*entity.User, error) {
	user, err := ur.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = ur.DB.DropRow(UserTableName, strconv.Itoa(id))
	if err != nil {
		return nil, ErrNoSuchUser
	}

	return user, nil
}

func (ur *UserRepoInMemDB) UpdateUser(ctx context.Context, id int, updated entity.User) (*entity.User, error) {
	user, err := ur.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	updated.ID = id
	updated.CreatedAt = user.CreatedAt
	updated.UpdatedAt = time.Now()

	err = ur.DB.AlterRow(UserTableName, strconv.Itoa(id), updated)
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
