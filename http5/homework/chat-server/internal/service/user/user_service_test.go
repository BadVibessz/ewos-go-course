package user

import (
	"context"
	sliceutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/slice"
	"math"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entity"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/mocks"

	repoerrors "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/repository"

	testingutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/testing"
)

func TestUserService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()
	now := time.Now()

	repoMock := mocks.NewMockUserRepo(ctrl)
	hasherMock := mocks.NewMockHasher(ctrl)

	service := New(repoMock, hasherMock)

	type inputArgs = entity.User
	type outputArg = *entity.User

	tests := []struct {
		name          string
		mockBehaviour func()
		input         inputArgs
		want          outputArg
		wantErr       bool
	}{
		{ // TODO: more test cases!
			name: "ok",
			mockBehaviour: func() {
				repoMock.
					EXPECT().
					CheckUniqueConstraints(ctx, "email@mail.com", "username").
					Return(nil)

				hasherMock.
					EXPECT().
					GenerateFromPassword([]byte("password"), 10).
					Return([]byte("hashed_password"), nil)

				repoMock.
					EXPECT().
					AddUser(
						ctx,
						entity.User{
							Email:          "email@mail.com",
							Username:       "username",
							HashedPassword: "hashed_password",
							CreatedAt:      now,
							UpdatedAt:      now,
						}).
					Return(
						&entity.User{
							ID:             1,
							Email:          "email@mail.com",
							Username:       "username",
							HashedPassword: "hashed_password",
							CreatedAt:      now,
							UpdatedAt:      now,
						},
						nil)
			},

			input: inputArgs{
				Email:          "email@mail.com",
				Username:       "username",
				HashedPassword: "password",
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
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehaviour()

			got, err := service.RegisterUser(ctx, test.input)

			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, testingutils.UsersEquals(*test.want, *got))
			}
		})
	}
}

func TestUserService_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()
	now := time.Now()

	repoMock := mocks.NewMockUserRepo(ctrl)
	hasherMock := mocks.NewMockHasher(ctrl)

	service := New(repoMock, hasherMock)

	type inputArgs = int
	type outputArg = *entity.User

	tests := []struct {
		name          string
		mockBehaviour func()
		input         inputArgs
		want          outputArg
		wantErr       bool
	}{
		{
			name: "ok, valid id",
			mockBehaviour: func() {
				repoMock.
					EXPECT().
					GetUserByID(ctx, 1).
					Return(
						&entity.User{
							ID:             1,
							Email:          "email@mail.com",
							Username:       "username",
							HashedPassword: "hashed_password",
							CreatedAt:      now,
							UpdatedAt:      now,
						},
						nil)
			},

			input: 1,
			want: &entity.User{
				ID:             1,
				Email:          "email@mail.com",
				Username:       "username",
				HashedPassword: "hashed_password",
				CreatedAt:      now,
				UpdatedAt:      now,
			},
		},
		{
			name: "err, invalid id (no such)",
			mockBehaviour: func() {
				repoMock.
					EXPECT().
					GetUserByID(ctx, 1).
					Return(nil, repoerrors.ErrNoSuchUser)
			},

			input:   1,
			wantErr: true,
		},
		{
			name: "err, invalid id (negative)",
			mockBehaviour: func() {
				repoMock.
					EXPECT().
					GetUserByID(ctx, -1).
					Return(nil, repoerrors.ErrNoSuchUser)
			},

			input:   -1,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehaviour()

			got, err := service.GetUserByID(ctx, test.input)

			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, testingutils.UsersEquals(*test.want, *got))
			}
		})
	}
}

func TestUserService_GetByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()
	now := time.Now()

	repoMock := mocks.NewMockUserRepo(ctrl)
	hasherMock := mocks.NewMockHasher(ctrl)

	service := New(repoMock, hasherMock)

	type inputArgs = string
	type outputArg = *entity.User

	tests := []struct {
		name          string
		mockBehaviour func()
		input         inputArgs
		want          outputArg
		wantErr       bool
	}{
		{
			name: "ok, valid email",
			mockBehaviour: func() {
				repoMock.
					EXPECT().
					GetUserByEmail(ctx, "email@mail.com").
					Return(
						&entity.User{
							ID:             1,
							Email:          "email@mail.com",
							Username:       "username",
							HashedPassword: "hashed_password",
							CreatedAt:      now,
							UpdatedAt:      now,
						},
						nil)
			},

			input: "email@mail.com",
			want: &entity.User{
				ID:             1,
				Email:          "email@mail.com",
				Username:       "username",
				HashedPassword: "hashed_password",
				CreatedAt:      now,
				UpdatedAt:      now,
			},
		},
		{
			name: "err, invalid email (no such)",
			mockBehaviour: func() {
				repoMock.
					EXPECT().
					GetUserByEmail(ctx, "email@mail.com").
					Return(nil, repoerrors.ErrNoSuchUser)
			},

			input:   "email@mail.com",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehaviour()

			got, err := service.GetUserByEmail(ctx, test.input)

			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, testingutils.UsersEquals(*test.want, *got))
			}
		})
	}
}

func TestUserService_GetByUsername(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()
	now := time.Now()

	repoMock := mocks.NewMockUserRepo(ctrl)
	hasherMock := mocks.NewMockHasher(ctrl)

	service := New(repoMock, hasherMock)

	type inputArgs = string
	type outputArg = *entity.User

	tests := []struct {
		name          string
		mockBehaviour func()
		input         inputArgs
		want          outputArg
		wantErr       bool
	}{
		{
			name: "ok, valid username",
			mockBehaviour: func() {
				repoMock.
					EXPECT().
					GetUserByUsername(ctx, "username").
					Return(
						&entity.User{
							ID:             1,
							Email:          "email@mail.com",
							Username:       "username",
							HashedPassword: "hashed_password",
							CreatedAt:      now,
							UpdatedAt:      now,
						},
						nil)
			},

			input: "username",
			want: &entity.User{
				ID:             1,
				Email:          "email@mail.com",
				Username:       "username",
				HashedPassword: "hashed_password",
				CreatedAt:      now,
				UpdatedAt:      now,
			},
		},
		{
			name: "err, invalid email (no such)",
			mockBehaviour: func() {
				repoMock.
					EXPECT().
					GetUserByUsername(ctx, "username").
					Return(nil, repoerrors.ErrNoSuchUser)
			},

			input:   "username",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehaviour()

			got, err := service.GetUserByUsername(ctx, test.input)

			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, testingutils.UsersEquals(*test.want, *got))
			}
		})
	}
}

func TestUserService_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()
	now := time.Now()

	repoMock := mocks.NewMockUserRepo(ctrl)
	hasherMock := mocks.NewMockHasher(ctrl)

	service := New(repoMock, hasherMock)

	type outputArg = []entity.User

	tests := []struct {
		name          string
		mockBehaviour func()
		offset        int
		limit         int
		want          outputArg
		wantErr       bool
	}{
		{
			name: "ok, no limit, no offset",
			mockBehaviour: func() {
				repoMock.
					EXPECT().
					GetAllUsers(ctx, 0, math.MaxInt64).
					Return(
						[]*entity.User{
							{
								ID:             1,
								Email:          "email@mail.com",
								Username:       "username",
								HashedPassword: "hashed_password",
								CreatedAt:      now,
								UpdatedAt:      now,
							},
							{
								ID:             2,
								Email:          "email2@mail.com",
								Username:       "username2",
								HashedPassword: "hashed_password",
								CreatedAt:      now,
								UpdatedAt:      now,
							},
							{
								ID:             3,
								Email:          "email3@mail.com",
								Username:       "username3",
								HashedPassword: "hashed_password",
								CreatedAt:      now,
								UpdatedAt:      now,
							},
						},
					)
			},

			offset: 0,
			limit:  math.MaxInt64,
			want: []entity.User{
				{
					ID:             1,
					Email:          "email@mail.com",
					Username:       "username",
					HashedPassword: "hashed_password",
					CreatedAt:      now,
					UpdatedAt:      now,
				},
				{
					ID:             2,
					Email:          "email2@mail.com",
					Username:       "username2",
					HashedPassword: "hashed_password",
					CreatedAt:      now,
					UpdatedAt:      now,
				},
				{
					ID:             3,
					Email:          "email3@mail.com",
					Username:       "username3",
					HashedPassword: "hashed_password",
					CreatedAt:      now,
					UpdatedAt:      now,
				},
			},
		},
		{
			name: "ok, no limit, offset 1",
			mockBehaviour: func() {
				repoMock.
					EXPECT().
					GetAllUsers(ctx, 1, math.MaxInt64).
					Return(
						[]*entity.User{
							{
								ID:             2,
								Email:          "email2@mail.com",
								Username:       "username2",
								HashedPassword: "hashed_password",
								CreatedAt:      now,
								UpdatedAt:      now,
							},
							{
								ID:             3,
								Email:          "email3@mail.com",
								Username:       "username3",
								HashedPassword: "hashed_password",
								CreatedAt:      now,
								UpdatedAt:      now,
							},
						},
					)
			},

			offset: 1,
			limit:  math.MaxInt64,
			want: []entity.User{
				{
					ID:             2,
					Email:          "email2@mail.com",
					Username:       "username2",
					HashedPassword: "hashed_password",
					CreatedAt:      now,
					UpdatedAt:      now,
				},
				{
					ID:             3,
					Email:          "email3@mail.com",
					Username:       "username3",
					HashedPassword: "hashed_password",
					CreatedAt:      now,
					UpdatedAt:      now,
				},
			},
		},
		{
			name: "ok, limit greater than data length, offset 1",
			mockBehaviour: func() {
				repoMock.
					EXPECT().
					GetAllUsers(ctx, 1, 10).
					Return(
						[]*entity.User{
							{
								ID:             2,
								Email:          "email2@mail.com",
								Username:       "username2",
								HashedPassword: "hashed_password",
								CreatedAt:      now,
								UpdatedAt:      now,
							},
							{
								ID:             3,
								Email:          "email3@mail.com",
								Username:       "username3",
								HashedPassword: "hashed_password",
								CreatedAt:      now,
								UpdatedAt:      now,
							},
						},
					)
			},

			offset: 1,
			limit:  10,
			want: []entity.User{
				{
					ID:             2,
					Email:          "email2@mail.com",
					Username:       "username2",
					HashedPassword: "hashed_password",
					CreatedAt:      now,
					UpdatedAt:      now,
				},
				{
					ID:             3,
					Email:          "email3@mail.com",
					Username:       "username3",
					HashedPassword: "hashed_password",
					CreatedAt:      now,
					UpdatedAt:      now,
				},
			},
		},
		{
			name: "ok, limit 1, offset 1", // TODO: NEGATIVE VALUE TEST CASE?
			mockBehaviour: func() {
				repoMock.
					EXPECT().
					GetAllUsers(ctx, 1, 1).
					Return(
						[]*entity.User{
							{
								ID:             2,
								Email:          "email2@mail.com",
								Username:       "username2",
								HashedPassword: "hashed_password",
								CreatedAt:      now,
								UpdatedAt:      now,
							},
						},
					)
			},

			offset: 1,
			limit:  1,
			want: []entity.User{
				{
					ID:             2,
					Email:          "email2@mail.com",
					Username:       "username2",
					HashedPassword: "hashed_password",
					CreatedAt:      now,
					UpdatedAt:      now,
				},
			},
		},
		{
			name: "ok, no limit, offset more than data length",
			mockBehaviour: func() {
				repoMock.
					EXPECT().
					GetAllUsers(ctx, 10, math.MaxInt64).
					Return(nil)
			},

			offset: 10,
			limit:  math.MaxInt64,
			want:   nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehaviour()

			got := service.GetAllUsers(ctx, test.offset, test.limit)

			assert.True(t, sliceutils.PointerAndValueSlicesEquals(got, test.want))
		})
	}
}

func TestUserService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()
	now := time.Now()

	repoMock := mocks.NewMockUserRepo(ctrl)
	hasherMock := mocks.NewMockHasher(ctrl)

	service := New(repoMock, hasherMock)

	type outputArg = *entity.User

	tests := []struct {
		name          string
		mockBehaviour func()
		id            int
		updateModel   entity.User
		want          outputArg
		wantErr       bool
	}{
		{
			name: "ok, valid id, valid update model",
			mockBehaviour: func() {
				repoMock.
					EXPECT().
					CheckUniqueConstraints(ctx, "newemail@mail.com", "newusername").
					Return(nil)

				repoMock.
					EXPECT().
					GetUserByID(ctx, 1).
					Return(
						&entity.User{
							ID:             1,
							Email:          "email@mail.com",
							Username:       "username",
							HashedPassword: "hashed_password",
							CreatedAt:      now,
							UpdatedAt:      now,
						},
						nil,
					)

				hasherMock.
					EXPECT().
					GenerateFromPassword([]byte("new_password"), 10).
					Return([]byte("new_hashed_password"), nil)

				repoMock.
					EXPECT().
					UpdateUser(
						ctx,
						1,
						entity.User{
							Email:          "newemail@mail.com",
							Username:       "newusername",
							HashedPassword: "new_hashed_password",
						},
					).
					Return(
						&entity.User{
							ID:             1,
							Email:          "newemail@mail.com",
							Username:       "newusername",
							HashedPassword: "new_hashed_password",
							CreatedAt:      now,
							UpdatedAt:      now,
						},
						nil)
			},
			id: 1,
			updateModel: entity.User{
				Email:          "newemail@mail.com",
				Username:       "newusername",
				HashedPassword: "new_password",
			},
			want: &entity.User{
				ID:             1,
				Email:          "newemail@mail.com",
				Username:       "newusername",
				HashedPassword: "new_hashed_password",
				CreatedAt:      now,
				UpdatedAt:      now,
			},
		},
		{
			name: "err, invalid id",
			mockBehaviour: func() {
				repoMock.
					EXPECT().
					CheckUniqueConstraints(ctx, "newemail@mail.com", "newusername").
					Return(nil)

				repoMock.
					EXPECT().
					GetUserByID(ctx, 1).
					Return(nil, repoerrors.ErrNoSuchUser)
			},
			id: 1,
			updateModel: entity.User{
				Email:          "newemail@mail.com",
				Username:       "newusername",
				HashedPassword: "new_hashed_password",
			},
			wantErr: true,
		},
		{
			name: "err, invalid update model (existing email)",
			mockBehaviour: func() {
				repoMock.
					EXPECT().
					CheckUniqueConstraints(ctx, "existingemail@mail.com", "").
					Return(repoerrors.ErrEmailExists)
			},
			id: 1,
			updateModel: entity.User{
				Email: "existingemail@mail.com",
			},
			wantErr: true,
		},
		{
			name: "err, invalid update model (existing username)",
			mockBehaviour: func() {
				repoMock.
					EXPECT().
					CheckUniqueConstraints(ctx, "", "existingusername").
					Return(repoerrors.ErrUsernameExists)
			},
			id: 1,
			updateModel: entity.User{
				Username: "existingusername",
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehaviour()

			got, err := service.UpdateUser(ctx, test.id, test.updateModel)

			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, testingutils.UsersEquals(*test.want, *got))
			}
		})
	}
}
