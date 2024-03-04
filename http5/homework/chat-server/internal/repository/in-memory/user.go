// nolint
package in_memory

import (
	"context"
	"errors"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/repository"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entity"

	inmemory "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/db/in-memory"
	sliceutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/slice"
)

type UserRepoInMemDB struct {
	mutex sync.RWMutex
	DB    inmemory.InMemoryDB
}

func NewInMemUserRepo(db inmemory.InMemoryDB) *UserRepoInMemDB {
	repo := UserRepoInMemDB{
		DB:    db,
		mutex: sync.RWMutex{},
	}

	_, err := repo.DB.GetTable(UserTableName)
	if errors.Is(err, inmemory.ErrNotExistedTable) {
		repo.DB.CreateTable(UserTableName)
	}

	return &repo
}

func (ur *UserRepoInMemDB) getAllUsers(_ context.Context, offset, limit int) []*entity.User {
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

func (ur *UserRepoInMemDB) GetAllUsers(ctx context.Context, offset, limit int) []*entity.User {
	ur.mutex.RLock()
	defer ur.mutex.RUnlock()

	return ur.getAllUsers(ctx, offset, limit)
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

func (ur *UserRepoInMemDB) getUserByID(_ context.Context, id int) (*entity.User, error) {
	row, err := ur.DB.GetRow(UserTableName, strconv.Itoa(id))
	if err != nil {
		return nil, repository.ErrNoSuchUser
	}

	user, ok := row.(entity.User)
	if !ok {
		return nil, repository.ErrNoSuchUser
	}

	return &user, nil
}

func (ur *UserRepoInMemDB) GetUserByID(ctx context.Context, id int) (*entity.User, error) {
	ur.mutex.RLock()
	defer ur.mutex.RUnlock()

	return ur.getUserByID(ctx, id)
}

func (ur *UserRepoInMemDB) getUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	users := ur.getAllUsers(ctx, 0, math.MaxInt64)
	if len(users) == 0 {
		return nil, repository.ErrNoSuchUser
	}

	filtered := sliceutils.Filter(users, func(u *entity.User) bool { return u.Email == email })
	if len(filtered) == 0 {
		return nil, repository.ErrNoSuchUser
	}

	user := filtered[0]

	return user, nil
}

func (ur *UserRepoInMemDB) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	ur.mutex.RLock()
	defer ur.mutex.RUnlock()

	return ur.getUserByEmail(ctx, email)
}

func (ur *UserRepoInMemDB) getUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	users := ur.getAllUsers(ctx, 0, math.MaxInt64)
	if len(users) == 0 {
		return nil, repository.ErrNoSuchUser
	}

	filtered := sliceutils.Filter(users, func(u *entity.User) bool { return u.Username == username })

	if len(filtered) == 0 {
		return nil, repository.ErrNoSuchUser
	}

	user := filtered[0]

	return user, nil
}

func (ur *UserRepoInMemDB) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	ur.mutex.RLock()
	defer ur.mutex.RUnlock()

	return ur.getUserByUsername(ctx, username)
}

func (ur *UserRepoInMemDB) DeleteUser(ctx context.Context, id int) (*entity.User, error) {
	ur.mutex.Lock()
	defer ur.mutex.Unlock()

	user, err := ur.getUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err = ur.DB.DropRow(UserTableName, strconv.Itoa(id)); err != nil {
		return nil, repository.ErrNoSuchUser
	}

	return user, nil
}

func (ur *UserRepoInMemDB) UpdateUser(ctx context.Context, id int, updated entity.User) (*entity.User, error) {
	ur.mutex.Lock()
	defer ur.mutex.Unlock()

	user, err := ur.getUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	updated.ID = id
	updated.CreatedAt = user.CreatedAt
	updated.UpdatedAt = time.Now()

	err = ur.DB.AlterRow(UserTableName, strconv.Itoa(id), updated)
	if err != nil {
		return nil, repository.ErrNoSuchUser
	}

	return user, nil
}

func (ur *UserRepoInMemDB) CheckUniqueConstraints(ctx context.Context, email, username string) error {
	ur.mutex.RLock()
	defer ur.mutex.RUnlock()

	got, err := ur.getUserByEmail(ctx, email)
	if got != nil || err == nil {
		return repository.ErrEmailExists
	}

	got, err = ur.getUserByUsername(ctx, username)
	if got != nil || err == nil {
		return repository.ErrUsernameExists
	}

	return nil
}
