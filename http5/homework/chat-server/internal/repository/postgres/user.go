package postgres

import (
	"context"
	"fmt"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/repository"
	"math"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entity"
)

type UserRepo struct {
	DB *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{
		DB: db,
	}
}

func (ur *UserRepo) GetAllUsers(ctx context.Context, offset, limit int) []*entity.User {
	var query string

	if limit == math.MaxInt64 {
		query = fmt.Sprintf("SELECT * FROM users ORDER BY created_at OFFSET %v", offset)
	} else {
		query = fmt.Sprintf("SELECT * FROM users ORDER BY created_at LIMIT %v OFFSET %v", limit, offset)
	}

	rows, err := ur.DB.QueryxContext(ctx, query)
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

	return users
}

func (ur *UserRepo) AddUser(ctx context.Context, user entity.User) (*entity.User, error) {
	now := time.Now()

	user.CreatedAt = now
	user.UpdatedAt = now

	result, err := ur.DB.NamedQueryContext(ctx,
		`INSERT INTO users (email, username, hashed_password, created_at, updated_at) 
VALUES (:email, :username, :hashed_password, :created_at, :updated_at) 
RETURNING id, email, username, hashed_password, created_at, updated_at`,
		&user)
	if err != nil {
		return nil, err
	}

	var usr entity.User

	if result.Next() {
		if err = result.StructScan(&usr); err != nil {
			return nil, err
		}
	}

	return &usr, nil
}

func (ur *UserRepo) getUserByArg(ctx context.Context, argName string, arg any) (*entity.User, error) {
	row := ur.DB.QueryRowxContext(ctx, fmt.Sprintf("SELECT * FROM users WHERE %v = '%v'", argName, arg))
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

func (ur *UserRepo) GetUserByID(ctx context.Context, id int) (*entity.User, error) {
	return ur.getUserByArg(ctx, "id", id)
}

func (ur *UserRepo) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return ur.getUserByArg(ctx, "email", email)
}

func (ur *UserRepo) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	return ur.getUserByArg(ctx, "username", username)
}

func (ur *UserRepo) DeleteUser(ctx context.Context, id int) (*entity.User, error) {
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

func (ur *UserRepo) UpdateUser(ctx context.Context, id int, updated entity.User) (*entity.User, error) {
	updated.UpdatedAt = time.Now()

	tx := ur.DB.MustBegin()

	query := "UPDATE users SET email=:email, username=:username, hashed_password=:hashed_password, updated_at=:updated_at" + fmt.Sprintf("WHERE :id = %v", id)

	_, err := tx.NamedExecContext(ctx, query, &updated)
	if err != nil {
		return nil, err
	}

	row := tx.QueryRowxContext(ctx, "SELECT * FROM users WHERE id = $1", id)
	if err = row.Err(); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	} // todo: maybe tx.Commit() or tx.Rollback() inside defer func?

	var user entity.User

	err = row.StructScan(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepo) CheckUniqueConstraints(ctx context.Context, email, username string) error {
	got, err := ur.GetUserByEmail(ctx, email)
	if got != nil || err == nil {
		return repository.ErrEmailExists
	}

	got, err = ur.GetUserByUsername(ctx, username)
	if got != nil || err == nil {
		return repository.ErrUsernameExists
	}

	return nil
}
