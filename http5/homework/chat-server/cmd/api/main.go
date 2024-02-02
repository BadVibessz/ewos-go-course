package main

import (
	"context"
	"errors"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/dto"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/repository"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/service"
	inmemory "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/db/in-memory"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/router"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/cmd/api/docs"
)

// @title Chat API
// @version 1.0
// @description API Server for Web Chat

// @BasePath /chat/api/v1

// @securityDefinitions.basic BasicAuth
// @in header
// @name Authorization

const dbSavePath = "http5/homework/chat-server/internal/db/db_state.json"

func initDB(ctx context.Context) (*inmemory.InMemDB, <-chan any) {
	var inMemDB *inmemory.InMemDB
	var savedChan <-chan any

	dbStateRestored := true

	jsonDb, err := os.ReadFile(dbSavePath)
	if err == nil {
		inMemDB, savedChan, err = inmemory.NewInMemDBFromJSON(ctx, string(jsonDb), dbSavePath)
		if err != nil {
			dbStateRestored = false
		}

	} else {
		dbStateRestored = false
	}

	if !dbStateRestored {
		inMemDB, savedChan = inmemory.NewInMemDB(ctx, dbSavePath)
	}

	return inMemDB, savedChan
}

func addDefaultDataInDB(ctx context.Context, us *service.UserService, ms *service.MessageService) {
	users := []dto.CreateUserDTO{
		{
			Username:       "test",
			Email:          "test@mail.ru",
			HashedPassword: "$2a$10$n1ZupQQL9NBnIDHShSIfwut3wf2cUMtsmzBo/7r29oRo4tYRrmoLS",
		},
		{
			Username:       "test2",
			Email:          "test2@mail.ru",
			HashedPassword: "$2a$10$O3bRPhNaWgVibnpkUFL.K.xXwmYnDKKMJ1Ak4iavFrSnn8wAsgYPW",
		},
		{
			Username:       "test3",
			Email:          "test3@mail.ru",
			HashedPassword: "$2a$10$lgQ9a71CwJQkAF1yUcKKl..RGDT4OaGRjyBAVFgGupkdMclmS7wMS",
		},
	}

	for _, user := range users {
		_, err := us.RegisterUser(ctx, user)
		if err != nil {
			return
		}
	}

	pubMessages := []dto.CreatePublicMessageDTO{
		{
			FromID:  1,
			Content: "Hello everyone, I'm Test!",
		},
		{
			FromID:  2,
			Content: "Hi Test, I'm Test2 ;)",
		},
		{
			FromID:  3,
			Content: "What's up! I'm Test3",
		},
	}

	for _, pubMsg := range pubMessages {
		_, err := ms.SendPublicMessage(ctx, pubMsg)
		if err != nil {
			return
		}
	}

	privMessages := []dto.CreatePrivateMessageDTO{
		{
			FromID:  1,
			ToID:    2,
			Content: "Excuse me, where am I?",
		},
		{
			FromID:  2,
			ToID:    1,
			Content: "Ohh.. You are being tested too!",
		},

		{
			FromID:  3,
			ToID:    2,
			Content: "Have something?",
		},
		{
			FromID:  2,
			ToID:    3,
			Content: "What??.. Get off me!",
		},
	}

	for _, privMsg := range privMessages {
		_, err := ms.SendPrivateMessage(ctx, privMsg)
		if err != nil {
			return
		}
	}

}

func main() {
	logger := logrus.New()

	ctx, cancel := context.WithCancel(context.Background())

	inMemDB, savedChan := initDB(ctx)

	userRepo := repository.NewInMemUserRepo(inMemDB)
	privateMsgRepo := repository.NewInMemPrivateMessageRepo(inMemDB)
	publicMsgRepo := repository.NewInMemPublicMessageRepo(inMemDB)

	userService := service.NewUserService(userRepo)
	messageService := service.NewMessageService(privateMsgRepo, publicMsgRepo, userRepo)
	authService := service.NewBasicAuthService(userRepo)

	addDefaultDataInDB(ctx, userService, messageService)

	userHandler := handler.NewUserHandler(userService, messageService, authService, logger)
	messageHandler := handler.NewMessageHandler(messageService, userService, authService, logger)

	routers := make(map[string]chi.Router)

	routers["/users"] = userHandler.Routes()
	routers["/messages"] = messageHandler.Routes()

	middlewares := []router.Middleware{
		middleware.Recoverer,
	}

	r := router.MakeRoutes("/chat/api/v1", routers, middlewares)

	server := http.Server{
		Addr:    ":5000",
		Handler: r,
	}

	logger.Infof("server started at %v", server.Addr)
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.WithError(err).Fatalf("server can't listen requests")
		}
	}()

	interrupt := make(chan os.Signal)

	signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)
	signal.Notify(interrupt, syscall.SIGINT)

	go func() {
		<-interrupt

		logger.Info("interrupt signal caught")
		logger.Info("chat api server shutting down")

		if err := server.Shutdown(ctx); err != nil {
			logger.WithError(err).Fatalf("can't close server listening on '%s'", server.Addr)
		}

		cancel()
	}()

	<-savedChan
}
