// nolint
package user

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entity"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/mapper"
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

type Middleware = func(http.Handler) http.Handler

type Handler struct {
	UserService    UserService
	MessageService MessageService
	Middlewares    []Middleware

	logger    *logrus.Logger
	validator *validator.Validate
}

func New(userService UserService,
	messageService MessageService,
	logger *logrus.Logger,
	validator *validator.Validate,
	middlewares ...Middleware,
) *Handler {
	return &Handler{
		UserService:    userService,
		MessageService: messageService,
		Middlewares:    middlewares,
		logger:         logger,
		validator:      validator,
	}
}

func (h *Handler) Routes() *chi.Mux {
	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		r.Use(h.Middlewares...)
		r.Get("/all", h.GetAll)
		r.Get("/messages", h.GetAllUsersThatSentMessage)
	})

	return router
}

// GetAll godoc
//
//	@Summary		Get all users
//	@Description	Get all users
//	@Security		BasicAuth
//	@Security		JWT
//	@Tags			User
//	@Produce		json
//	@Success		200	{object}	[]response.GetUserResponse
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
//	@Security		JWT
//	@Tags			User
//	@Produce		json
//	@Success		200	{object}	[]response.GetUserResponse
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
