package auth

import (
	"context"
	"errors"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entity"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/mocks"
	testingutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/testing"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	repoerrors "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/repository"
)

func TestAuthService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()
	now := time.Now()

	userRepoMock := mocks.NewMockUserRepo(ctrl)
	hasherMock := mocks.NewMockAuthHasher(ctrl)

	service := New(userRepoMock, hasherMock)

	type inputArgs = entity.User
	type outputArg = *entity.User

	tests := []struct {
		name          string
		mockBehaviour func()
		username      string
		password      string
		want          outputArg
		wantErr       bool
	}{
		{
			name: "ok, valid username and password",
			mockBehaviour: func() {
				userRepoMock.
					EXPECT().
					GetUserByUsername(ctx, "username").
					Return(&entity.User{
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
					CompareHashAndPassword([]byte("hashed_password"), []byte("password")).
					Return(nil)
			},

			username: "username",
			password: "password",
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
			name: "err, invalid username, valid password",
			mockBehaviour: func() {
				userRepoMock.
					EXPECT().
					GetUserByUsername(ctx, "invalid_username").
					Return(nil, repoerrors.ErrNoSuchUser)
			},

			username: "invalid_username",
			password: "password",
			wantErr:  true,
		},
		{
			name: "err, valid username, invalid password",
			mockBehaviour: func() {
				userRepoMock.
					EXPECT().
					GetUserByUsername(ctx, "username").
					Return(&entity.User{
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
					CompareHashAndPassword([]byte("hashed_password"), []byte("invalid_password")).
					Return(errors.New("hashes are not equal"))
			},

			username: "username",
			password: "invalid_password",
			wantErr:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehaviour()

			got, err := service.Login(ctx, test.username, test.password)

			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, testingutils.UsersEquals(*test.want, *got))
			}
		})
	}
}
