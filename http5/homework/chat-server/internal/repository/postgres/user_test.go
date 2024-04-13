package postgres

import (
	"context"
	"errors"
	testingutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/pkg/utils/testing"
	"math"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	sqlxmock "github.com/zhashkevych/go-sqlxmock"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entity"

	sliceutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/slice"
)

func TestUserRepo_AddUser(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

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
				rows := sqlxmock.
					NewRows([]string{"id", "username", "email", "hashed_password", "created_at", "updated_at"}).
					AddRow(1, "username", "email@mail.com", "hashed_password", now, now)

				mock.ExpectQuery("INSERT INTO users").
					WithArgs("email@mail.com", "username", "hashed_password", testingutils.AnyTime{}, testingutils.AnyTime{}).
					WillReturnRows(rows)
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
				mock.ExpectQuery("INSERT INTO users").
					WithArgs("", "", "", testingutils.AnyTime{}, testingutils.AnyTime{}).
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
				assert.True(t, testingutils.UsersEquals(*test.want, *got))
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

	defer db.Close()

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
				assert.True(t, sliceutils.PointerAndValueSlicesEquals(got, test.want))
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepo_GetByID(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	repo := NewUserRepo(db)

	type outputArg = entity.User
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
					NewRows([]string{"id", "username", "email", "hashed_password", "created_at", "updated_at"}).
					AddRow(1, "username", "email@mail.com", "hashed_password", time.Time{}, time.Time{})

				mock.ExpectQuery("SELECT * FROM users WHERE id = 1").
					WillReturnRows(rows)
			},
			input: 1,
			want: entity.User{
				ID:             1,
				Username:       "username",
				Email:          "email@mail.com",
				HashedPassword: "hashed_password",
				CreatedAt:      time.Time{},
				UpdatedAt:      time.Time{},
			},
		},
		{
			name: "err, invalid id",
			mockBehaviour: func() {
				rows := sqlxmock.
					NewRows([]string{"id", "username", "email", "hashed_password", "created_at", "updated_at"})

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM users WHERE $1 = $2`)).
					WithArgs("id", 2).
					WillReturnRows(rows)
			},
			input:   2,
			wantErr: true,
		},
	}

	ctx := context.Background()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehaviour()

			got, err := repo.GetUserByID(ctx, test.input)

			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want, *got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepo_GetByUsername(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	repo := NewUserRepo(db)

	type outputArg = entity.User
	type inputArg = string

	tests := []struct {
		name          string
		mockBehaviour func()
		input         inputArg
		want          outputArg
		wantErr       bool
	}{
		{
			name: "ok, valid username", // TODO: NOT FUCKING WORKING
			mockBehaviour: func() {
				rows := sqlxmock.
					NewRows([]string{"id", "username", "email", "hashed_password", "created_at", "updated_at"}).
					AddRow(1, "username", "email@mail.com", "hashed_password", time.Time{}, time.Time{})

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM users WHERE $1 = $2`)).
					WithArgs("username", "username").
					WillReturnRows(rows)
			},
			input: "username",
			want: entity.User{
				ID:             1,
				Username:       "username",
				Email:          "email@mail.com",
				HashedPassword: "hashed_password",
				CreatedAt:      time.Time{},
				UpdatedAt:      time.Time{},
			},
		},
		{
			name: "err, invalid username",
			mockBehaviour: func() {
				rows := sqlxmock.
					NewRows([]string{"id", "username", "email", "hashed_password", "created_at", "updated_at"})

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM users WHERE $1 = $2`)).
					WithArgs("username", "not_presented").
					WillReturnRows(rows)
			},
			input:   "not_presented",
			wantErr: true,
		},
	}

	ctx := context.Background()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehaviour()

			got, err := repo.GetUserByUsername(ctx, test.input)

			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want, *got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepo_GetByEmail(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	repo := NewUserRepo(db)

	type outputArg = entity.User
	type inputArg = string

	tests := []struct {
		name          string
		mockBehaviour func()
		input         inputArg
		want          outputArg
		wantErr       bool
	}{
		{
			name: "ok, valid email",
			mockBehaviour: func() {
				rows := sqlxmock.
					NewRows([]string{"id", "username", "email", "hashed_password", "created_at", "updated_at"}).
					AddRow(1, "username", "email@mail.com", "hashed_password", time.Time{}, time.Time{})

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM users WHERE $1 = $2`)).
					WithArgs("email", "email@mail.com").
					WillReturnRows(rows)
			},
			input: "email@mail.com",
			want: entity.User{
				ID:             1,
				Username:       "username",
				Email:          "email@mail.com",
				HashedPassword: "hashed_password",
				CreatedAt:      time.Time{},
				UpdatedAt:      time.Time{},
			},
		},
		{
			name: "err, invalid email",
			mockBehaviour: func() {
				rows := sqlxmock.
					NewRows([]string{"id", "username", "email", "hashed_password", "created_at", "updated_at"})

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM users WHERE $1 = $2`)).
					WithArgs("email", "not_presented@mail.com").
					WillReturnRows(rows)
			},
			input:   "not_presented@mail.com",
			wantErr: true,
		},
	}

	ctx := context.Background()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehaviour()

			got, err := repo.GetUserByEmail(ctx, test.input)

			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want, *got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepo_Delete(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	repo := NewUserRepo(db)

	type outputArg = entity.User
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
					NewRows([]string{"id", "username", "email", "hashed_password", "created_at", "updated_at"}).
					AddRow(1, "username", "email@mail.com", "hashed_password", time.Time{}, time.Time{})

				mock.ExpectQuery(regexp.QuoteMeta(`DELETE FROM users WHERE id = $1`)).
					WithArgs(1).
					WillReturnRows(rows)
			},
			input: 1,
			want: entity.User{
				ID:             1,
				Username:       "username",
				Email:          "email@mail.com",
				HashedPassword: "hashed_password",
				CreatedAt:      time.Time{},
				UpdatedAt:      time.Time{},
			},
		},
		{
			name: "err, invalid id",
			mockBehaviour: func() {
				rows := sqlxmock.
					NewRows([]string{"id", "username", "email", "hashed_password", "created_at", "updated_at"})

				mock.ExpectQuery(regexp.QuoteMeta(`DELETE FROM users WHERE id = $1`)).
					WithArgs(1).
					WillReturnRows(rows)
			},
			input:   1,
			wantErr: true,
		},
	}

	ctx := context.Background()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehaviour()

			got, err := repo.DeleteUser(ctx, test.input)

			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want, *got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

// TODO: TEST repo.UpdateUser
