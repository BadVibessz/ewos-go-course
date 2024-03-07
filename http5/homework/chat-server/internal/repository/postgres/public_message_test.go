package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	sqlxmock "github.com/zhashkevych/go-sqlxmock"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entity"
)

func publicMessagesEqual(msg1, msg2 entity.PublicMessage) bool {
	return msg1.ID == msg2.ID &&
		msg1.FromUsername == msg2.FromUsername &&
		msg1.Content == msg2.Content &&
		timesAlmostEqual(msg1.SentAt, msg2.SentAt) &&
		timesAlmostEqual(msg1.EditedAt, msg2.EditedAt)
}

func TestPublicMessageRepo_AddMessage(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	repo := NewPublicMessageRepo(db)

	type inputArgs = entity.PublicMessage
	type outputArg = *entity.PublicMessage

	now := time.Now()

	tests := []struct {
		name          string
		mockBehaviour func()
		input         inputArgs
		want          outputArg
		wantErr       bool
	}{
		{
			name: "ok",
			mockBehaviour: func() {
				result := sqlxmock.NewResult(1, 1)

				mock.ExpectExec("INSERT INTO public_message").
					WithArgs("from_username", "content", AnyTime{}, AnyTime{}).
					WillReturnResult(result)
			},

			input: inputArgs{
				FromUsername: "from_username",
				Content:      "content",
			},
			want: &inputArgs{
				ID:           1,
				FromUsername: "from_username",
				Content:      "content",
				SentAt:       now,
				EditedAt:     now,
			},
		},
		{
			name: "empty fields",
			mockBehaviour: func() {
				mock.ExpectExec("INSERT INTO public_message").
					WithArgs("", "", AnyTime{}, AnyTime{}).
					WillReturnError(errors.New("not null constraint not satisfied"))
			},

			input: inputArgs{
				FromUsername: "",
				Content:      "",
			},

			wantErr: true,
		},
	}

	ctx := context.Background()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehaviour()

			got, err := repo.AddPublicMessage(ctx, test.input)

			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, publicMessagesEqual(*test.want, *got))
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

//func TestPublicMessageRepo_GetAll(t *testing.T) {
//	db, mock, err := sqlxmock.Newx()
//	if err != nil {
//		t.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
//	}
//
//	defer db.Close()
//
//	repo := NewPublicMessageRepo(db)
//
//	type outputArg = []entity.PublicMessage
//
//	tests := []struct {
//		name          string
//		mockBehaviour func()
//		limit         int
//		offset        int
//		want          outputArg
//		wantErr       bool
//	}{
//		{
//			name: "ok, no limit, no offset",
//			mockBehaviour: func() {
//				rows := sqlxmock.
//					NewRows([]string{"id", "from_username", "content", "sent_at", "edited_at"}).
//					AddRow(1, "username", "content", time.Time{}, time.Time{}).
//					AddRow(2, "username2", "content", time.Time{}, time.Time{}).
//					AddRow(3, "username3", "content", time.Time{}, time.Time{}).
//					mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM users ORDER BY created_at OFFSET 0`)).WillReturnRows(rows)
//			},
//
//			limit:  math.MaxInt64,
//			offset: 0,
//			want: []entity.User{
//				{
//					ID:             1,
//					Username:       "username",
//					Email:          "email@mail.com",
//					HashedPassword: "hashed_password",
//					CreatedAt:      time.Time{},
//					UpdatedAt:      time.Time{},
//				},
//				{
//					ID:             2,
//					Username:       "username2",
//					Email:          "email2@mail.com",
//					HashedPassword: "hashed_password",
//					CreatedAt:      time.Time{},
//					UpdatedAt:      time.Time{},
//				},
//				{
//					ID:             3,
//					Username:       "username3",
//					Email:          "email3@mail.com",
//					HashedPassword: "hashed_password",
//					CreatedAt:      time.Time{},
//					UpdatedAt:      time.Time{},
//				},
//			},
//		},
//		{
//			name: "ok, no limit, offset 1",
//			mockBehaviour: func() {
//				rows := sqlxmock.
//					NewRows([]string{"id", "username", "email", "hashed_password", "created_at", "updated_at"}).
//					AddRow(2, "username2", "email2@mail.com", "hashed_password", time.Time{}, time.Time{}).
//					AddRow(3, "username3", "email3@mail.com", "hashed_password", time.Time{}, time.Time{})
//
//				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM users ORDER BY created_at OFFSET 1`)).WillReturnRows(rows)
//			},
//
//			limit:  math.MaxInt64,
//			offset: 1,
//			want: []entity.User{
//				{
//					ID:             2,
//					Username:       "username2",
//					Email:          "email2@mail.com",
//					HashedPassword: "hashed_password",
//					CreatedAt:      time.Time{},
//					UpdatedAt:      time.Time{},
//				},
//				{
//					ID:             3,
//					Username:       "username3",
//					Email:          "email3@mail.com",
//					HashedPassword: "hashed_password",
//					CreatedAt:      time.Time{},
//					UpdatedAt:      time.Time{},
//				},
//			},
//		},
//		{
//			name: "ok, limit 1, offset 1",
//			mockBehaviour: func() {
//				rows := sqlxmock.
//					NewRows([]string{"id", "username", "email", "hashed_password", "created_at", "updated_at"}).
//					AddRow(2, "username2", "email2@mail.com", "hashed_password", time.Time{}, time.Time{})
//
//				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM users ORDER BY created_at LIMIT 1 OFFSET 1`)).WillReturnRows(rows)
//			},
//
//			limit:  1,
//			offset: 1,
//			want: []entity.User{
//				{
//					ID:             2,
//					Username:       "username2",
//					Email:          "email2@mail.com",
//					HashedPassword: "hashed_password",
//					CreatedAt:      time.Time{},
//					UpdatedAt:      time.Time{},
//				},
//			},
//		},
//		{
//			name: "ok, limit -1, offset -1",
//			mockBehaviour: func() {
//				rows := sqlxmock.NewRows([]string{"id", "username", "email", "hashed_password", "created_at", "updated_at"})
//
//				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM users ORDER BY created_at LIMIT -1 OFFSET -1`)).
//					WillReturnRows(rows)
//			},
//
//			limit:  -1,
//			offset: -1,
//			want:   nil,
//		},
//	}
//
//	ctx := context.Background()
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			test.mockBehaviour()
//
//			got := repo.GetAllUsers(ctx, test.offset, test.limit)
//
//			if test.wantErr {
//				assert.Error(t, err)
//			} else {
//				assert.NoError(t, err)
//				assert.True(t, sliceutils.PointerAndValueSlicesEqual(got, test.want))
//			}
//
//			assert.NoError(t, mock.ExpectationsWereMet())
//		})
//	}
//}
