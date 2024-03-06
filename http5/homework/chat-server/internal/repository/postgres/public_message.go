package postgres

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entity"
)

type PublicMessageRepo struct {
	DB *sqlx.DB
}

func NewPublicMessageRepo(db *sqlx.DB) *PublicMessageRepo {
	return &PublicMessageRepo{
		DB: db,
	}
}

func (pr *PublicMessageRepo) AddPublicMessage(ctx context.Context, msg entity.PublicMessage) (*entity.PublicMessage, error) {
	now := time.Now()

	msg.SentAt = now
	msg.EditedAt = now

	result, err := pr.DB.NamedExecContext(ctx,
		"INSERT INTO public_message (from_username, content, sent_at, edited_at) VALUES (:from_username, :content, :sent_at, :edited_at)",
		&msg)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	msg.ID = int(id)

	return &msg, nil
}

func (pr *PublicMessageRepo) GetAllPublicMessages(ctx context.Context, offset, limit int) []*entity.PublicMessage {
	rows, err := pr.DB.QueryxContext(ctx, "SELECT * FROM public_message ORDER BY sent_at LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil
	}

	var users []*entity.PublicMessage

	for rows.Next() {
		var msg entity.PublicMessage

		err = rows.StructScan(&msg)
		if err != nil {
			return nil
		}

		users = append(users, &msg)
	}

	return users
}

func (pr *PublicMessageRepo) GetPublicMessage(ctx context.Context, id int) (*entity.PublicMessage, error) {
	row := pr.DB.QueryRowxContext(ctx, "SELECT * FROM public_message WHERE id = $1", id)
	if err := row.Err(); err != nil {
		return nil, err
	}

	var msg entity.PublicMessage

	err := row.StructScan(&msg)
	if err != nil {
		return nil, err
	}

	return &msg, nil
}
