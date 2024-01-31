package handler

import (
	"context"
	"fmt"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/dto"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/requset"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/middleware"
	handlerutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/pkg/utils/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type UserHandler struct {
	UserService middleware.UserService
	logger      *logrus.Logger
	validator   *validator.Validate
}

func NewUserHandler(us middleware.UserService, logger *logrus.Logger) *UserHandler {
	return &UserHandler{
		UserService: us,
		logger:      logger,
		validator:   validator.New(validator.WithRequiredStructEnabled()),
	}
}

// Routes todo: maybe here accept context and then use callback with context?
func (uh *UserHandler) Routes() chi.Router {
	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		r.Post("/register", uh.Register)
	})

	router.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(uh.UserService))
		r.Get("/all", uh.GetAll)
	})

	return router
}

func (uh *UserHandler) Register(rw http.ResponseWriter, req *http.Request) { // TODO: PANIC IF TRYING REGISTER 2nd USER
	registerReq := requset.RegisterRequest{}
	err := render.DecodeJSON(req.Body, &registerReq)
	if err != nil {
		logMsg := fmt.Sprintf("error occured decoding request body to RegisterRequest struct: %s", err)
		respMsg := fmt.Sprintf("invalid registration data provided")

		handlerutils.WriteResponseAndLogError(rw, uh.logger, http.StatusBadRequest, logMsg, respMsg)

		return
	}

	err = uh.validator.Struct(registerReq)
	if err != nil {
		logMsg := fmt.Sprintf("error occurred validating RegisterRequest struct: %s", err)
		respMsg := fmt.Sprintf("invalid registration data provided")

		handlerutils.WriteResponseAndLogError(rw, uh.logger, http.StatusBadRequest, logMsg, respMsg)

		return
	}

	// todo: understand what is bcrypt.cost!
	hash, err := bcrypt.GenerateFromPassword([]byte(registerReq.Password), bcrypt.DefaultCost)

	createModel := dto.CreateUserDTO{
		Email:          registerReq.Email,
		Username:       registerReq.Username,
		HashedPassword: string(hash),
	}

	// todo: understand how to manage ctx here
	ctx := context.Background()

	user, err := uh.UserService.RegisterUser(ctx, createModel)
	if err != nil {
		handlerutils.WriteResponseAndLogError(rw, uh.logger, http.StatusBadRequest, // todo: specify what user did wrong
			"", "invalid registration data provided")

		return
	}

	render.JSON(rw, req, user)
	rw.WriteHeader(http.StatusCreated)
}

func (uh *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	users := uh.UserService.GetAllUsers(ctx)

	render.JSON(w, r, users)

	w.WriteHeader(http.StatusOK)
}
