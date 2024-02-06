package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/model"
	"github.com/sirupsen/logrus"

	handlerutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/handler"
)

type Middleware = func(http.Handler) http.Handler

var (
	ErrNullUsername = errors.New("empty username provided")
	ErrNullPassword = errors.New("empty password provided")
)

type AuthService interface {
	Login(ctx context.Context, username, password string) (*model.User, error)
}

func AuthMiddleware(authService AuthService, logger *logrus.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			username, pass, ok := req.BasicAuth()
			if !ok {
				rw.WriteHeader(http.StatusUnauthorized)
				return
			}

			if username == "" {
				logMsg := fmt.Sprintf("error occurred while logging user: %s", ErrNullUsername)
				respMsg := fmt.Sprintf("error occurred while logging user: %s", ErrNullUsername)

				handlerutils.WriteErrResponseAndLog(rw, logger, http.StatusUnauthorized, logMsg, respMsg)

				return
			}

			if pass == "" {
				logMsg := fmt.Sprintf("error occurred while logging user: %s", ErrNullPassword)
				respMsg := fmt.Sprintf("error occurred while logging user: %s", ErrNullPassword)

				handlerutils.WriteErrResponseAndLog(rw, logger, http.StatusUnauthorized, logMsg, respMsg)

				return
			}

			user, err := authService.Login(req.Context(), username, pass)
			if err != nil {
				logMsg := fmt.Sprintf("error occurred while logging user: %s", err)
				respMsg := fmt.Sprintf("error occurred while logging user: %s", err)

				handlerutils.WriteErrResponseAndLog(rw, logger, http.StatusUnauthorized, logMsg, respMsg)

				return
			}

			req.Header.Set("id", strconv.Itoa(user.ID))

			next.ServeHTTP(rw, req)
		})
	}
}
