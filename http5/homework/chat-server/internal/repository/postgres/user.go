// nolint
package postgres

import (
	"context"
	"fmt"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entity"
	"github.com/jmoiron/sqlx"
	"sync"
	"time"
)

type UserPostgresRepo struct {
	mutex sync.RWMutex
	DB    *sqlx.DB
}

func New(db *sqlx.DB) *UserPostgresRepo {
	return &UserPostgresRepo{
		DB:    db,
		mutex: sync.RWMutex{},
	}
}

// todo: test
func (ur *UserPostgresRepo) GetAllUsers(ctx context.Context, offset, limit int) []*entity.User {
	// todo: listen for context?

	rows, err := ur.DB.QueryxContext(ctx, "SELECT * FROM users ORDER BY created_at LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil // todo: return err?
	}

	var users []*entity.User

	for rows.Next() {
		var user entity.User

		err = rows.StructScan(&user)
		if err != nil {
			return nil // todo: return err?
		}

		users = append(users, &user)
	}

	// map []User -> []*User
	//return sliceutils.Map(users, func(usr entity.User) *entity.User { return &usr })
	return users
}

func (ur *UserPostgresRepo) AddUser(ctx context.Context, user entity.User) (*entity.User, error) {
	now := time.Now()

	user.CreatedAt = now
	user.UpdatedAt = now

	result, err := ur.DB.NamedExecContext(ctx,
		"INSERT INTO users (email, username, hashed_password, created_at, updated_at) VALUES (:email, :username, :hashed_password, :created_at, :updated_at)",
		&user)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	user.ID = int(id)

	return &user, nil
}

func (ur *UserPostgresRepo) getUserByArg(ctx context.Context, argName string, arg any) (*entity.User, error) {
	row := ur.DB.QueryRowxContext(ctx, "SELECT * FROM users WHERE $1 = $2", argName, arg)
	if err := row.Err(); err != nil {
		return nil, err
	}

	var user entity.User

	err := row.StructScan(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserPostgresRepo) GetUserByID(ctx context.Context, id int) (*entity.User, error) {
	return ur.getUserByArg(ctx, "id", id)
}

func (ur *UserPostgresRepo) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return ur.getUserByArg(ctx, "email", email)
}

func (ur *UserPostgresRepo) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	return ur.getUserByArg(ctx, "username", username)
}

func (ur *UserPostgresRepo) DeleteUser(ctx context.Context, id int) (*entity.User, error) {
	row := ur.DB.QueryRowxContext(ctx, "DELETE FROM users WHERE id = $1", id)
	if err := row.Err(); err != nil {
		return nil, err
	}

	var user entity.User

	err := row.StructScan(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserPostgresRepo) UpdateUser(ctx context.Context, id int, updated entity.User) (*entity.User, error) {
	updated.UpdatedAt = time.Now()

	tx := ur.DB.MustBegin()

	query := "UPDATE users SET email=:email, username=:username, hashed_password=:hashed_password, updated_at=:updated_at" + fmt.Sprintf("WHERE :id = %v", id)

	_, err := tx.NamedExecContext(ctx, query, &updated)
	if err != nil {
		return nil, err
	}

	// todo: THIS PIECE OF CODE TO PRIVATE METHOD AS THIS PIECE MEETS OFTEN
	row := tx.QueryRowxContext(ctx, "SELECT * FROM users WHERE id = $1", id)
	if err = row.Err(); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	var user entity.User

	err = row.StructScan(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
