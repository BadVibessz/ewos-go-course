package main

import (
	"context"
	"errors"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/repository"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/service"
	inmemory "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/db/in-memory"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/router"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := logrus.New()

	ctx, cancel := context.WithCancel(context.Background())

	inMemDB := inmemory.NewInMemDB(ctx, "db_state.json")

	userRepo := repository.NewInMemUserRepo(inMemDB)
	privateMsgRepo := repository.NewInMemPrivateMessageRepo(inMemDB)
	publicMsgRepo := repository.NewInMemPublicMessageRepo(inMemDB)

	userService := service.NewUserService(userRepo)
	messageService := service.NewMessageService(privateMsgRepo, publicMsgRepo, userRepo)

	userHandler := handler.NewUserHandler(userService, logger)
	messageHandler := handler.NewMessageHandler(messageService, userService, logger)

	r := router.MakeRoutes("/chat/api/v1", []chi.Router{
		userHandler.Routes(),
		messageHandler.Routes(),
	})

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

	waitChan := make(chan any)

	go func() {
		<-interrupt

		logger.Info("interrupt signal caught")
		logger.Info("chat api server shutting down")

		if err := server.Shutdown(ctx); err != nil {
			logger.WithError(err).Fatalf("can't close server listening on '%s'", server.Addr)
		}

		cancel()
		waitChan <- any(0)
	}()

	<-waitChan
}
