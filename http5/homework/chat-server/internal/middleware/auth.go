package middleware

import (
	"context"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/dto"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/model"
	"net/http"
)

type Middleware = func(http.Handler) http.Handler

type AuthService interface {
	Login(ctx context.Context, cred dto.LoginUserDTO) (*model.User, error)
}

func AuthMiddleware(authService AuthService) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			username, pass, ok := req.BasicAuth()
			if !ok {
				rw.WriteHeader(http.StatusUnauthorized)
				return
			}

			cred := dto.LoginUserDTO{
				Username: username,
				Password: pass,
			}

			user, err := authService.Login(req.Context(), cred)
			if err != nil {
				rw.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(req.Context(), "user", *user)
			next.ServeHTTP(rw, req.WithContext(ctx))
		})
	}
}
