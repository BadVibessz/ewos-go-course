package handler

import (
	"context"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/dto"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/requset"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/middleware"
	"github.com/go-chi/chi"
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

	router.Route("/users", func(r chi.Router) {
		r.Post("/register", uh.Register)
	})

	router.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(uh.UserService))
		r.Get("/users/all", uh.GetAll)
	})

	return router
}

func (uh *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	registerReq := requset.RegisterRequest{}
	err := render.DecodeJSON(r.Body, &registerReq)
	if err != nil {
		uh.logger.Errorf("error occured decoding request body to RegisterRequest struct: %s", err)
		return
	}

	err = uh.validator.Struct(registerReq)
	if err != nil {
		uh.logger.Errorf("error occured validating RegisterRequest struct: %s", err)
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
		return
	}

	render.JSON(w, r, user)
	w.WriteHeader(http.StatusCreated)
}

func (uh *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	users := uh.UserService.GetAllUsers(ctx)

	render.JSON(w, r, users)

	w.WriteHeader(http.StatusOK)
}
