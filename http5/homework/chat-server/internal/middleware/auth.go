package middleware

import (
	"context"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/dto"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/model"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type Middleware = func(http.Handler) http.Handler

// UserService TODO: this interface should be in user.handler but there's circular import, how to resolve?
type UserService interface {
	RegisterUser(ctx context.Context, user dto.CreateUserDTO) (*model.User, error)
	GetUserByID(ctx context.Context, id int) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetAllUsers(ctx context.Context) []*model.User
	UpdateUser(ctx context.Context, id int, updateModel dto.UpdateUserDTO) (*model.User, error)
	DeleteUser(ctx context.Context, id int) (*model.User, error)
}

func AuthMiddleware(userService UserService) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			username, pass, ok := req.BasicAuth()
			if !ok {
				rw.WriteHeader(http.StatusUnauthorized)
				return
			}

			// todo: understand context management
			user, err := userService.GetUserByUsername(req.Context(), username)
			if err != nil {
				rw.WriteHeader(http.StatusUnauthorized)
				return
			}

			err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(pass))
			if err != nil {
				rw.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(req.Context(), "user", user)
			next.ServeHTTP(rw, req.WithContext(ctx))
		})
	}
}
