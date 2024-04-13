package postgres

import (
	"context"
	"errors"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entity"
	testingutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/pkg/utils/testing"
	"github.com/stretchr/testify/assert"
	"math"
	"regexp"
	"testing"
	"time"

	sqlxmock "github.com/zhashkevych/go-sqlxmock"

	sliceutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/slice"
)

func TestPrivateMessageRepo_AddMessage(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	repo := NewPrivateMessageRepo(db)

	type inputArgs = entity.PrivateMessage
	type outputArg = *entity.PrivateMessage

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
				rows := sqlxmock.NewRows([]string{"id", "from_username", "to_username", "content", "sent_at", "edited_at"}).
					AddRow(1, "from_username", "to_username", "content", now, now)

				mock.ExpectQuery("INSERT INTO private_message").
					WithArgs("from_username", "to_username", "content", testingutils.AnyTime{}, testingutils.AnyTime{}).
					WillReturnRows(rows)
			},

			input: inputArgs{
				FromUsername: "from_username",
				ToUsername:   "to_username",
				Content:      "content",
			},
			want: &inputArgs{
				ID:           1,
				FromUsername: "from_username",
				ToUsername:   "to_username",
				Content:      "content",
				SentAt:       now,
				EditedAt:     now,
			},
		},
		{
			name: "empty fields",
			mockBehaviour: func() {
				mock.ExpectQuery("INSERT INTO private_message").
					WithArgs("", "", "", testingutils.AnyTime{}, testingutils.AnyTime{}).
					WillReturnError(errors.New("not null constraint not satisfied"))
			},

			input: inputArgs{
				FromUsername: "",
				ToUsername:   "",
				Content:      "",
			},

			wantErr: true,
		},
	}

	ctx := context.Background()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehaviour()

			got, err := repo.AddPrivateMessage(ctx, test.input)

			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, testingutils.PrivateMessagesEquals(*test.want, *got))
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPrivateMessageRepo_GetAll(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	repo := NewPrivateMessageRepo(db)

	type outputArg = []entity.PrivateMessage

	tests := []struct {
		name          string
		mockBehaviour func()
		limit         int
		offset        int
		want          outputArg
		wantErr       bool
	}{
		{
			name: "ok, no limit, no offset",
			mockBehaviour: func() {
				rows := sqlxmock.
					NewRows([]string{"id", "from_username", "to_username", "content", "sent_at", "edited_at"}).
					AddRow(1, "from_username", "to_username", "content", time.Time{}, time.Time{}).
					AddRow(2, "from_username2", "to_username2", "content", time.Time{}, time.Time{}).
					AddRow(3, "from_username3", "to_username3", "content", time.Time{}, time.Time{})

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM private_message ORDER BY sent_at OFFSET 0`)).WillReturnRows(rows)
			},

			limit:  math.MaxInt64,
			offset: 0,
			want: []entity.PrivateMessage{
				{
					ID:           1,
					FromUsername: "from_username",
					ToUsername:   "to_username",
					Content:      "content",
					SentAt:       time.Time{},
					EditedAt:     time.Time{},
				},
				{
					ID:           2,
					FromUsername: "from_username2",
					ToUsername:   "to_username2",
					Content:      "content",
					SentAt:       time.Time{},
					EditedAt:     time.Time{},
				},
				{
					ID:           3,
					FromUsername: "from_username3",
					ToUsername:   "to_username3",
					Content:      "content",
					SentAt:       time.Time{},
					EditedAt:     time.Time{},
				},
			},
		},
		{
			name: "ok, no limit, offset 1",
			mockBehaviour: func() {
				rows := sqlxmock.
					NewRows([]string{"id", "from_username", "to_username", "content", "sent_at", "edited_at"}).
					AddRow(2, "from_username2", "to_username2", "content", time.Time{}, time.Time{}).
					AddRow(3, "from_username3", "to_username3", "content", time.Time{}, time.Time{})

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM private_message ORDER BY sent_at OFFSET 1`)).WillReturnRows(rows)
			},

			limit:  math.MaxInt64,
			offset: 1,
			want: []entity.PrivateMessage{
				{
					ID:           2,
					FromUsername: "from_username2",
					ToUsername:   "to_username2",
					Content:      "content",
					SentAt:       time.Time{},
					EditedAt:     time.Time{},
				},
				{
					ID:           3,
					FromUsername: "from_username3",
					ToUsername:   "to_username3",
					Content:      "content",
					SentAt:       time.Time{},
					EditedAt:     time.Time{},
				},
			},
		},
		{
			name: "ok, limit 1, offset 1",
			mockBehaviour: func() {
				rows := sqlxmock.
					NewRows([]string{"id", "from_username", "to_username", "content", "sent_at", "edited_at"}).
					AddRow(2, "from_username2", "to_username2", "content", time.Time{}, time.Time{})

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM private_message ORDER BY sent_at LIMIT 1 OFFSET 1`)).WillReturnRows(rows)
			},

			limit:  1,
			offset: 1,
			want: []entity.PrivateMessage{
				{
					ID:           2,
					FromUsername: "from_username2",
					ToUsername:   "to_username2",
					Content:      "content",
					SentAt:       time.Time{},
					EditedAt:     time.Time{},
				},
			},
		},
		{
			name: "ok, limit -1, offset -1",
			mockBehaviour: func() {
				rows := sqlxmock.NewRows([]string{"id", "username", "email", "hashed_password", "created_at", "updated_at"})

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM private_message ORDER BY sent_at LIMIT -1 OFFSET -1`)).
					WillReturnRows(rows)
			},

			limit:  -1,
			offset: -1,
			want:   nil,
		},
	}

	ctx := context.Background()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehaviour()

			got := repo.GetAllPrivateMessages(ctx, test.offset, test.limit)

			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, sliceutils.PointerAndValueSlicesEquals(got, test.want))
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPrivateMessageRepo_Get(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	repo := NewPrivateMessageRepo(db)

	type outputArg = entity.PrivateMessage
	type inputArg = int

	tests := []struct {
		name          string
		mockBehaviour func()
		input         inputArg
		want          outputArg
		wantErr       bool
	}{
		{
			name: "ok, valid id",
			mockBehaviour: func() {
				rows := sqlxmock.
					NewRows([]string{"id", "from_username", "to_username", "content", "sent_at", "edited_at"}).
					AddRow(1, "from_username", "to_username", "content", time.Time{}, time.Time{})

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM private_message WHERE id = $1`)).
					WithArgs(1).
					WillReturnRows(rows)
			},

			input: 1,
			want: entity.PrivateMessage{
				ID:           1,
				FromUsername: "from_username",
				ToUsername:   "to_username",
				Content:      "content",
				SentAt:       time.Time{},
				EditedAt:     time.Time{},
			},
		},
		{
			name: "err, invalid id (no such id)",
			mockBehaviour: func() {
				rows := sqlxmock.NewRows([]string{"id", "from_username", "to_username", "content", "sent_at", "edited_at"})

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM private_message WHERE id = $1`)).
					WithArgs(1).
					WillReturnRows(rows)
			},

			input:   1,
			wantErr: true,
		},
		{
			name: "err, invalid id (negative value)",
			mockBehaviour: func() {
				rows := sqlxmock.NewRows([]string{"id", "from_username", "to_username", "content", "sent_at", "edited_at"})

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM private_message WHERE id = $1`)).
					WithArgs(-1).
					WillReturnRows(rows)
			},

			input:   -1,
			wantErr: true,
		},
	}

	ctx := context.Background()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehaviour()

			got, err := repo.GetPrivateMessage(ctx, test.input)

			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, *got, test.want)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
