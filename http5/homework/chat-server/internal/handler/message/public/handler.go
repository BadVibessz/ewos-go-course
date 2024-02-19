// nolint
package public

import (
	"context"
	"fmt"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entity"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/mapper"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/middleware"
	handlerinternalutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/pkg/utils/handler"
	"github.com/go-playground/validator/v10"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/request"

	handlerutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/handler"
	sliceutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/slice"
)

type PublicMessageService interface {
	SendPublicMessage(ctx context.Context, fromID int, content string) (*entity.PublicMessage, error)
	GetPublicMessage(ctx context.Context, id int) (*entity.PublicMessage, error)
	GetAllPublicMessages(ctx context.Context, offset, limit int) []*entity.PublicMessage
}

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
	Login(ctx context.Context, loginReq request.LoginRequest) (*entity.User, error)
}

type Handler struct {
	MessageService PublicMessageService
	UserService    UserService
	AuthService    AuthService
	logger         *logrus.Logger
	valid          *validator.Validate
}

func New(publicMessageService PublicMessageService, userService UserService, authService AuthService,
	logger *logrus.Logger, valid *validator.Validate) *Handler {
	return &Handler{
		MessageService: publicMessageService,
		UserService:    userService,
		AuthService:    authService,
		logger:         logger,
		valid:          valid,
	}
}

func (h *Handler) Routes() *chi.Mux {
	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(h.AuthService, h.logger, h.valid))

		r.Get("/", h.GetAllPublicMessages)
		r.Post("/", h.SendPublicMessage)
	})

	return router
}

// GetAllPublicMessages godoc
//
//	@Summary		Get all public messages
//	@Description	Get all public messages that were sent to chat
//	@Security		BasicAuth
//	@Tags			Message
//	@Produce		json
//	@Param			offset	query		int	true	"Offset"
//	@Param			limit	query		int	true	"Limit"
//	@Success		200		{object}	[]response.PublicMessageResponse
//	@Failure		401		{string}	Unauthorized
//	@Router			/api/v1/messages/public [get]
func (h *Handler) GetAllPublicMessages(rw http.ResponseWriter, req *http.Request) {
	paginationOpts := handlerinternalutils.GetPaginationOptsFromQuery(req, handler.DefaultOffset, handler.DefaultLimit)

	if err := paginationOpts.Validate(h.valid); err != nil {
		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, "", err.Error())

		return
	}

	messages := h.MessageService.GetAllPublicMessages(req.Context(), paginationOpts.Offset, paginationOpts.Limit)

	render.JSON(rw, req, sliceutils.Map(messages, mapper.MapPublicMessageToResponse))
	rw.WriteHeader(http.StatusOK)
}

// SendPublicMessage godoc
//
//	@Summary		Send public message to chat
//	@Description	Send public message to chat
//	@Security		BasicAuth
//	@Tags			Message
//	@Accept			json
//	@Produce		json
//	@Param			input	body		request.SendPublicMessageRequest	true	"public message schema"
//	@Success		200		{object}	[]response.PublicMessageResponse
//	@Failure		401		{string}	Unauthorized
//	@Failure		500		{string}	internal	error
//	@Router			/api/v1/messages/public [post]
func (h *Handler) SendPublicMessage(rw http.ResponseWriter, req *http.Request) {
	id, err := handlerutils.GetIntHeaderByKey(req, "id")
	if err != nil {
		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusUnauthorized, "", err.Error())
		return
	}

	var pubMsgReq request.SendPublicMessageRequest

	if err = render.DecodeJSON(req.Body, &pubMsgReq); err != nil {
		logMsg := fmt.Sprintf("error occurred validating PublicMessageRequest struct: %s", err)
		respMsg := fmt.Sprintf("invalid message provided: %s", err)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, logMsg, respMsg)

		return
	}

	pubMsgReq.FromID = id

	if err = pubMsgReq.Validate(h.valid); err != nil {
		logMsg := fmt.Sprintf("error occurred validating PublicMessageRequest struct: %s", err)
		respMsg := fmt.Sprintf("invalid message provided: %s", err)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, logMsg, respMsg)

		return
	}

	message, err := h.MessageService.SendPublicMessage(req.Context(), pubMsgReq.FromID, pubMsgReq.Content)
	if err != nil {
		logMsg := fmt.Sprintf("error occurred saving public message: %s", err)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusInternalServerError, logMsg, "")

		return
	}

	render.JSON(rw, req, mapper.MapPublicMessageToResponse(message))
	rw.WriteHeader(http.StatusCreated)
}
