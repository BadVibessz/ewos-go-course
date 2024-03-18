package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entity"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/middleware/mapper"
	handlerutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/handler"
	jwtutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/jwt"
)

type Handler = func(http.Handler) http.Handler

type AuthService interface {
	Login(ctx context.Context, username, password string) (*entity.User, error)
}

func BasicAuthMiddleware(authService AuthService, logger *logrus.Logger, valid *validator.Validate) Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			loginReq, err := mapper.MapBasicAuthToLoginRequest(req.BasicAuth())
			if err != nil {
				msg := fmt.Sprintf("error occurred while logging user: %v", err)

				handlerutils.WriteErrResponseAndLog(rw, logger, http.StatusBadRequest, msg, msg)
				return
			}

			if err = loginReq.Validate(valid); err != nil {
				msg := fmt.Sprintf("error occurred while logging user: %v", err)

				handlerutils.WriteErrResponseAndLog(rw, logger, http.StatusBadRequest, msg, msg)
				return
			}

			user, err := authService.Login(req.Context(), loginReq.Username, loginReq.Password)
			if err != nil {
				msg := fmt.Sprintf("error occurred while logging user: %v", err)

				handlerutils.WriteErrResponseAndLog(rw, logger, http.StatusUnauthorized, msg, msg)
				return
			}

			req.Header.Set("id", strconv.Itoa(user.ID))
			req.Header.Set("username", user.Username)

			next.ServeHTTP(rw, req)
		})
	}
}

func JWTAuthMiddleware(secret string, logger *logrus.Logger) Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			authHeader := req.Header.Get("Authorization")
			if authHeader == "" {
				msg := "authorization header is empty"

				handlerutils.WriteErrResponseAndLog(rw, logger, http.StatusUnauthorized, msg, msg)
				return
			}

			token := authHeader[len("Bearer "):] // TODO: JWT ACCESS AND REFRESH TOKEN

			payload, err := jwtutils.ValidateToken(token, secret)
			if err != nil {
				msg := fmt.Sprintf("error occurred validating token: %v", err)

				handlerutils.WriteErrResponseAndLog(rw, logger, http.StatusUnauthorized, msg, msg)
				return
			}

			// todo: to private func
			idAny, exists := payload["id"]
			if !exists {
				msg := "invalid payload: not contains id"

				handlerutils.WriteErrResponseAndLog(rw, logger, http.StatusUnauthorized, msg, msg)
				return
			}

			id, ok := idAny.(float64)
			if !ok {
				msg := "cannot parse id from payload to float64"

				handlerutils.WriteErrResponseAndLog(rw, logger, http.StatusUnauthorized, msg, msg)
				return
			}

			usernameAny, exists := payload["username"]
			if !exists {
				msg := "invalid payload: not contains username"

				handlerutils.WriteErrResponseAndLog(rw, logger, http.StatusUnauthorized, msg, msg)
				return
			}

			username, ok := usernameAny.(string)
			if !ok {
				msg := "cannot parse username from payload to string"

				handlerutils.WriteErrResponseAndLog(rw, logger, http.StatusUnauthorized, msg, msg)
				return
			}

			req.Header.Set("id", strconv.Itoa(int(id)))
			req.Header.Set("username", username)

			next.ServeHTTP(rw, req)
		})
	}
}
