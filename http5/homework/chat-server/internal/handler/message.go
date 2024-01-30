package handler

import (
	"context"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/dto"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/middleware"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/model"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
	"net/http"
)

type MessageService interface {
	SendPrivateMessage(ctx context.Context, createModel dto.CreatePrivateMessageDTO) (*model.PrivateMessage, error)
	SendPublicMessage(ctx context.Context, createModel dto.CreatePublicMessageDTO) (*model.PublicMessage, error)

	GetPrivateMessage(ctx context.Context, id int) (*model.PrivateMessage, error)
	GetPublicMessage(ctx context.Context, id int) (*model.PublicMessage, error)

	GetAllPrivateMessages(ctx context.Context) []*model.PrivateMessage
	GetAllPublicMessages(ctx context.Context) []*model.PublicMessage

	UpdatePrivateMessage(ctx context.Context, id int, newContent string) (*model.PrivateMessage, error)
	UpdatePublicMessage(ctx context.Context, id int, newContent string) (*model.PublicMessage, error)

	DeletePrivateMessage(ctx context.Context, id int) (*model.PrivateMessage, error)
	DeletePublicMessage(ctx context.Context, id int) (*model.PublicMessage, error)
}

type MessageHandler struct {
	MessageService MessageService
	UserService    middleware.UserService
	logger         *logrus.Logger
}

func NewMessageHandler(ms MessageService, us middleware.UserService, logger *logrus.Logger) *MessageHandler {
	return &MessageHandler{
		MessageService: ms,
		UserService:    us,
		logger:         logger,
	}
}

func (mh *MessageHandler) Routes() chi.Router {
	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(mh.UserService))
		r.Get("/messages/public", mh.GetAllPublicMessages)
		r.Post("/messages/public", mh.SendPublicMessage)
	})

	return router
}

func (mh *MessageHandler) GetAllPublicMessages(w http.ResponseWriter, req *http.Request) {
	msgs := mh.MessageService.GetAllPublicMessages(req.Context())

	render.JSON(w, req, msgs)

	w.WriteHeader(http.StatusOK)
}

func (mh *MessageHandler) SendPublicMessage(rw http.ResponseWriter, req *http.Request) {
	user, ok := req.Context().Value("user").(model.User)
	if !ok { // unauthorized
		rw.WriteHeader(http.StatusUnauthorized)
	}

	var content string

	err := render.DecodeJSON(req.Body, &content)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if content == "" {
		rw.WriteHeader(http.StatusBadRequest)
		_, err := rw.Write([]byte("Message cannot be empty!"))
		if err != nil {
			mh.logger.Errorf("error occured writing response: %s", err)
		}

		return
	}

	msg := dto.CreatePublicMessageDTO{ // todo: maybe dto should use mode.User struct?
		FromID:  user.ID,
		Content: content,
	}

	message, err := mh.MessageService.SendPublicMessage(req.Context(), msg)
	if err != nil {
		mh.logger.Errorf("error occurred saving public message: %s", err)
		rw.WriteHeader(http.StatusInternalServerError) // todo: correct status code?
		return
	}

	render.JSON(rw, req, message)
	rw.WriteHeader(http.StatusCreated)
}
