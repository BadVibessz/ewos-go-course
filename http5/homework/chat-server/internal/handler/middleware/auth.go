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
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/request"

	handlerutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/handler"
)

type Handler = func(http.Handler) http.Handler

type AuthService interface {
	Login(ctx context.Context, loginReq request.LoginRequest) (*entity.User, error)
}

func BasicAuthMiddleware(authService AuthService, logger *logrus.Logger, valid *validator.Validate) Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			loginReq, err := mapper.MapBasicAuthToLoginRequest(req.BasicAuth())
			if err != nil {
				logMsg := fmt.Sprintf("error occurred while logging user: %v", err)
				respMsg := fmt.Sprintf("error occurred while logging user: %v", err)

				handlerutils.WriteErrResponseAndLog(rw, logger, http.StatusBadRequest, logMsg, respMsg)

				return
			}

			if err = loginReq.Validate(valid); err != nil {
				logMsg := fmt.Sprintf("error occurred while logging user: %v", err)
				respMsg := fmt.Sprintf("error occurred while logging user: %v", err)

				handlerutils.WriteErrResponseAndLog(rw, logger, http.StatusBadRequest, logMsg, respMsg)

				return
			}

			user, err := authService.Login(req.Context(), *loginReq)
			if err != nil {
				logMsg := fmt.Sprintf("error occurred while logging user: %v", err)
				respMsg := fmt.Sprintf("error occurred while logging user: %v", err)

				handlerutils.WriteErrResponseAndLog(rw, logger, http.StatusUnauthorized, logMsg, respMsg)

				return
			}

			req.Header.Set("id", strconv.Itoa(user.ID))

			next.ServeHTTP(rw, req)
		})
	}
}

func JWTAuthMiddleware(authService AuthService, logger *logrus.Logger, valid *validator.Validate) Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

			// todo: check jwt token passed in auth header
		})
	}
}
