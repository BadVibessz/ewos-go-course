package postgres

import (
	"context"
	"database/sql/driver"
	"errors"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entity"
	sliceutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/slice"
	"math"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	sqlxmock "github.com/zhashkevych/go-sqlxmock"
)

func timesAlmostEqual(tim1, tim2 time.Time) bool {
	return tim1.Sub(tim2) <= 1*time.Second
}

func usersEqual(usr1, usr2 entity.User) bool {
	return usr1.ID == usr2.ID &&
		usr1.Username == usr2.Username &&
		usr1.Email == usr2.Email &&
		timesAlmostEqual(usr1.CreatedAt, usr2.CreatedAt) &&
		timesAlmostEqual(usr1.UpdatedAt, usr2.UpdatedAt)

}

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestUserRepo_AddUser(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
	}

	defer func() { // TODO: WHY ERROR
		if err = db.Close(); err != nil {
			t.Fatalf("an error '%v' was not expected when closing a stub database connection", err)
		}
	}()

	repo := NewUserRepo(db)

	type inputArgs = entity.User
	type outputArg = *entity.User

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

				mock.ExpectExec("INSERT INTO users").
					WithArgs("email@mail.com", "username", "hashed_password", AnyTime{}, AnyTime{}).
					WillReturnResult(result)
			},

			input: inputArgs{
				Email:          "email@mail.com",
				Username:       "username",
				HashedPassword: "hashed_password",
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			want: &inputArgs{
				ID:             1,
				Email:          "email@mail.com",
				Username:       "username",
				HashedPassword: "hashed_password",
				CreatedAt:      now,
				UpdatedAt:      now,
			},
		},
		{
			name: "empty fields",
			mockBehaviour: func() {
				mock.ExpectExec("INSERT INTO users").
					WithArgs("", "", "", AnyTime{}, AnyTime{}).
					WillReturnError(errors.New("not null constraint not satisfied"))
			},

			input: inputArgs{
				Email:          "",
				Username:       "",
				HashedPassword: "",
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			want: &inputArgs{
				ID:             1,
				Email:          "",
				Username:       "",
				HashedPassword: "",
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			wantErr: true,
		},
		// todo: add test cases?
	}

	ctx := context.Background()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehaviour()

			got, err := repo.AddUser(ctx, test.input)

			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, usersEqual(*test.want, *got))
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepo_GetAll(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
	}

	defer func() { // TODO: WHY ERROR
		if err = db.Close(); err != nil {
			t.Fatalf("an error '%v' was not expected when closing a stub database connection", err)
		}
	}()

	repo := NewUserRepo(db)

	type outputArg = []entity.User

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
					NewRows([]string{"id", "username", "email", "hashed_password", "created_at", "updated_at"}).
					AddRow(1, "username", "email@mail.com", "hashed_password", time.Time{}, time.Time{}).
					AddRow(2, "username2", "email2@mail.com", "hashed_password", time.Time{}, time.Time{}).
					AddRow(3, "username3", "email3@mail.com", "hashed_password", time.Time{}, time.Time{})

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM users ORDER BY created_at OFFSET 0`)).WillReturnRows(rows)
			},

			limit:  math.MaxInt64,
			offset: 0,
			want: []entity.User{
				{
					ID:             1,
					Username:       "username",
					Email:          "email@mail.com",
					HashedPassword: "hashed_password",
					CreatedAt:      time.Time{},
					UpdatedAt:      time.Time{},
				},
				{
					ID:             2,
					Username:       "username2",
					Email:          "email2@mail.com",
					HashedPassword: "hashed_password",
					CreatedAt:      time.Time{},
					UpdatedAt:      time.Time{},
				},
				{
					ID:             3,
					Username:       "username3",
					Email:          "email3@mail.com",
					HashedPassword: "hashed_password",
					CreatedAt:      time.Time{},
					UpdatedAt:      time.Time{},
				},
			},
		},
		{
			name: "ok, no limit, offset 1",
			mockBehaviour: func() {
				rows := sqlxmock.
					NewRows([]string{"id", "username", "email", "hashed_password", "created_at", "updated_at"}).
					AddRow(2, "username2", "email2@mail.com", "hashed_password", time.Time{}, time.Time{}).
					AddRow(3, "username3", "email3@mail.com", "hashed_password", time.Time{}, time.Time{})

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM users ORDER BY created_at OFFSET 1`)).WillReturnRows(rows)
			},

			limit:  math.MaxInt64,
			offset: 1,
			want: []entity.User{
				{
					ID:             2,
					Username:       "username2",
					Email:          "email2@mail.com",
					HashedPassword: "hashed_password",
					CreatedAt:      time.Time{},
					UpdatedAt:      time.Time{},
				},
				{
					ID:             3,
					Username:       "username3",
					Email:          "email3@mail.com",
					HashedPassword: "hashed_password",
					CreatedAt:      time.Time{},
					UpdatedAt:      time.Time{},
				},
			},
		},
		{
			name: "ok, limit 1, offset 1",
			mockBehaviour: func() {
				rows := sqlxmock.
					NewRows([]string{"id", "username", "email", "hashed_password", "created_at", "updated_at"}).
					AddRow(2, "username2", "email2@mail.com", "hashed_password", time.Time{}, time.Time{})

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM users ORDER BY created_at LIMIT 1 OFFSET 1`)).WillReturnRows(rows)
			},

			limit:  1,
			offset: 1,
			want: []entity.User{
				{
					ID:             2,
					Username:       "username2",
					Email:          "email2@mail.com",
					HashedPassword: "hashed_password",
					CreatedAt:      time.Time{},
					UpdatedAt:      time.Time{},
				},
			},
		},
		{
			name: "ok, limit -1, offset -1",
			mockBehaviour: func() {
				rows := sqlxmock.NewRows([]string{"id", "username", "email", "hashed_password", "created_at", "updated_at"})

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM users ORDER BY created_at LIMIT -1 OFFSET -1`)).
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

			got := repo.GetAllUsers(ctx, test.offset, test.limit)

			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, sliceutils.PointerAndValueSlicesEqual(got, test.want))
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
