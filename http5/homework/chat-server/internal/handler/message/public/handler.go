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
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/request"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/model"

	handlerutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/handler"
	sliceutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/slice"
)

type PublicMessageService interface {
	SendPublicMessage(ctx context.Context, createModel entity.PublicMessage) (*model.PublicMessage, error)
	GetPublicMessage(ctx context.Context, id int) (*model.PublicMessage, error)
	GetAllPublicMessages(ctx context.Context, paginationOpts request.PaginationOptions) []*model.PublicMessage
}

type UserService interface {
	RegisterUser(ctx context.Context, user entity.User) (*model.User, error)
	GetUserByID(ctx context.Context, id int) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetAllUsers(ctx context.Context, paginationOpts request.PaginationOptions) []*model.User
	UpdateUser(ctx context.Context, id int, updateModel entity.User) (*model.User, error)
	DeleteUser(ctx context.Context, id int) (*model.User, error)
}

type AuthService interface {
	Login(ctx context.Context, loginReq request.LoginRequest) (*model.User, error)
}

type Handler struct {
	MessageService PublicMessageService
	UserService    UserService
	AuthService    AuthService
	logger         *logrus.Logger
}

func New(publicMessageService PublicMessageService, userService UserService, authService AuthService, logger *logrus.Logger) *Handler {
	return &Handler{
		MessageService: publicMessageService,
		UserService:    userService,
		AuthService:    authService,
		logger:         logger,
	}
}

func (h *Handler) Routes() *chi.Mux {
	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(h.AuthService, h.logger))

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

	if err := paginationOpts.Validate(); err != nil {
		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, "", err.Error())

		return
	}

	messages := h.MessageService.GetAllPublicMessages(req.Context(), paginationOpts)

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

	if err = pubMsgReq.Validate(); err != nil {
		logMsg := fmt.Sprintf("error occurred validating PublicMessageRequest struct: %s", err)
		respMsg := fmt.Sprintf("invalid message provided: %s", err)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, logMsg, respMsg)

		return
	}

	msg := mapper.MapPublicMessageRequestToEntity(&pubMsgReq)

	message, err := h.MessageService.SendPublicMessage(req.Context(), msg)
	if err != nil {
		logMsg := fmt.Sprintf("error occurred saving public message: %s", err)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusInternalServerError, logMsg, "")

		return
	}

	render.JSON(rw, req, mapper.MapPublicMessageToResponse(message))
	rw.WriteHeader(http.StatusCreated)
}
