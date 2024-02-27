// nolint
package user

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entity"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/mapper"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/middleware"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/request"

	handlerinternalutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/pkg/utils/handler"
	handlerutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/handler"
	sliceutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/slice"
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

type MessageService interface {
	GetAllPrivateMessages(ctx context.Context, toID int, offset, limit int) []*entity.PrivateMessage
	GetAllUsersThatSentMessage(ctx context.Context, toID int, offset, limit int) []*entity.User
}

type AuthService interface {
	Login(ctx context.Context, loginReq request.LoginRequest) (*entity.User, error)
}

type Handler struct {
	UserService    UserService
	MessageService MessageService
	AuthService    AuthService
	logger         *logrus.Logger
	validator      *validator.Validate
}

func New(userService UserService,
	messageService MessageService,
	authService AuthService,
	logger *logrus.Logger,
	validator *validator.Validate,
) *Handler {
	return &Handler{
		UserService:    userService,
		MessageService: messageService,
		AuthService:    authService,
		logger:         logger,
		validator:      validator,
	}
}

func (h *Handler) Routes() *chi.Mux {
	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		r.Post("/register", h.Register)
	})

	router.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(h.AuthService, h.logger, h.validator))
		r.Get("/all", h.GetAll)
		r.Get("/messages", h.GetAllUsersThatSentMessage)
	})

	return router
}

// Register godoc
//
//	@Summary		Register new user
//	@Description	to register new user
//	@Tags			User
//	@Accept			json
//	@Produce		plain
//	@Param			input	body		request.RegisterRequest	true	"registration info"
//	@Success		200		{object}	response.UserResponse
//	@Failure		400		{string}	invalid		registration	data	provided
//	@Failure		500		{string}	internal	error
//	@Router			/api/v1/users/register [post]
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

// GetAll godoc
//
//	@Summary		Get all users
//	@Description	Get all users
//	@Security		BasicAuth
//	@Tags			User
//	@Produce		json
//	@Success		200	{object}	[]response.UserResponse
//	@Failure		401	{string}	Unauthorized
//	@Router			/api/v1/users/all [get]
func (h *Handler) GetAll(rw http.ResponseWriter, req *http.Request) {
	paginationOpts := handlerinternalutils.GetPaginationOptsFromQuery(req, handler.DefaultOffset, handler.DefaultLimit)

	if err := paginationOpts.Validate(h.validator); err != nil {
		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, "", err.Error())

		return
	}

	users := h.UserService.GetAllUsers(req.Context(), paginationOpts.Offset, paginationOpts.Limit)

	render.JSON(rw, req, sliceutils.Map(users, mapper.MapUserToUserResponse))
	rw.WriteHeader(http.StatusOK)
}

// GetAllUsersThatSentMessage godoc
//
//	@Summary		Get all users that sent message to current user
//	@Description	Get all users that sent message to current user
//	@Security		BasicAuth
//	@Tags			User
//	@Produce		json
//	@Success		200	{object}	[]response.UserResponse
//	@Failure		401	{string}	Unauthorized
//	@Router			/api/v1/users/messages [get]
func (h *Handler) GetAllUsersThatSentMessage(rw http.ResponseWriter, req *http.Request) {
	id, err := handlerutils.GetIntHeaderByKey(req, "id")
	if err != nil {
		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusUnauthorized, "", err.Error())
		return
	}

	paginateOpts := request.GetUnlimitedPaginationOptions()

	users := h.MessageService.GetAllUsersThatSentMessage(req.Context(), id, paginateOpts.Offset, paginateOpts.Limit)

	render.JSON(rw, req, sliceutils.Map(users, mapper.MapUserToUserResponse))
	rw.WriteHeader(http.StatusOK)
}
