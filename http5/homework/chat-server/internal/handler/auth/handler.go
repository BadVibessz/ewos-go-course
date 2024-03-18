package auth

import (
	"context"
	"fmt"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/config"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entity"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/mapper"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/request"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/response"

	handlerutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/handler"
	jwtutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/jwt"
)

type UserService interface {
	RegisterUser(ctx context.Context, user entity.User) (*entity.User, error)
	GetUserByID(ctx context.Context, id int) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
	GetAllUsers(ctx context.Context, offset, limit int) []*entity.User
	UpdateUser(ctx context.Context, id int, updateModel entity.User) (*entity.User, error)
	DeleteUser(ctx context.Context, id int) (*entity.User, error)
}

type AuthService interface {
	Login(ctx context.Context, username, password string) (*entity.User, error)
}

type Middleware = func(http.Handler) http.Handler

type Handler struct {
	UserService UserService
	AuthService AuthService
	Middlewares []Middleware

	JwtConfig config.Jwt

	logger    *logrus.Logger
	validator *validator.Validate
}

func New(userService UserService,
	authService AuthService,
	jwtConfig config.Jwt,
	logger *logrus.Logger,
	validator *validator.Validate,
	middlewares ...Middleware,
) *Handler {
	return &Handler{
		UserService: userService,
		AuthService: authService,
		JwtConfig:   jwtConfig,
		Middlewares: middlewares,
		logger:      logger,
		validator:   validator,
	}
}

func (h *Handler) Routes() *chi.Mux {
	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		r.Use(h.Middlewares...)
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
	})

	return router
}

// Register godoc
//
//	@Summary		Register new user
//	@Description	to register new user
//	@Tags			Auth
//	@Accept			json
//	@Produce		plain
//	@Param			input	body		request.RegisterRequest	true	"registration info"
//	@Success		200		{object}	response.GetUserResponse
//	@Failure		400		{string}	invalid		registration	data	provided
//	@Failure		500		{string}	internal	error
//	@Router			/api/v1/auth/register [post]
func (h *Handler) Register(rw http.ResponseWriter, req *http.Request) {
	var registerReq request.RegisterRequest

	if err := render.DecodeJSON(req.Body, &registerReq); err != nil {
		logMsg := fmt.Sprintf("error occurred decoding request body to RegisterRequest struct: %v", err)
		respMsg := fmt.Sprintf("invalid registration data provided: %v", err)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, logMsg, respMsg)

		return
	}

	if err := registerReq.Validate(h.validator); err != nil {
		logMsg := fmt.Sprintf("error occurred validating RegisterRequest struct: %v", err)
		respMsg := fmt.Sprintf("invalid registration data provided: %v", err)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, logMsg, respMsg)

		return
	}

	user, err := h.UserService.RegisterUser(req.Context(), mapper.MapRegisterRequestToUserEntity(&registerReq))
	if err != nil {
		msg := fmt.Sprintf("error occurred registrating user: %v", err)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusInternalServerError, msg, msg)

		return
	}

	render.JSON(rw, req, mapper.MapUserToUserResponse(user))
	rw.WriteHeader(http.StatusCreated)
}

// Login godoc
//
//	@Summary		Login user
//	@Description	login user via JWT
//	@Tags			Auth
//	@Accept			json
//	@Produce		plain
//	@Param			input	body		request.LoginRequest	true	"login info"
//	@Success		200		{object}	response.LoginResponse
//	@Failure		400		{string}	invalid		login	data	provided
//	@Failure		500		{string}	internal	error
//	@Router			/api/v1/auth/login [post]
func (h *Handler) Login(rw http.ResponseWriter, req *http.Request) {
	var loginReq request.LoginRequest

	if err := render.DecodeJSON(req.Body, &loginReq); err != nil {
		logMsg := fmt.Sprintf("error occurred decoding request body to LoginRequest struct: %v", err)
		respMsg := fmt.Sprintf("invalid login data provided: %v", err)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, logMsg, respMsg)
		return
	}

	if err := loginReq.Validate(h.validator); err != nil {
		logMsg := fmt.Sprintf("error occurred validating LoginRequest struct: %v", err)
		respMsg := fmt.Sprintf("invalid login data provided: %v", err)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, logMsg, respMsg)
		return
	}

	user, err := h.AuthService.Login(req.Context(), loginReq.Username, loginReq.Password)
	if err != nil {
		msg := fmt.Sprintf("error occurred while user login: %v", err)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, msg, msg)
		return
	}

	// construct jwt token
	payload := map[string]any{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	}

	token, err := jwtutils.CreateJWT(payload, jwt.SigningMethodHS256, h.JwtConfig.Secret)
	if err != nil {
		msg := fmt.Sprintf("error occurred signing jwt token: %v", err)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusInternalServerError, msg, msg)
		return
	}

	// todo: mapper?
	resp := response.LoginResponse{
		Token: token,
	}

	render.JSON(rw, req, resp)
	rw.WriteHeader(http.StatusOK)
}
