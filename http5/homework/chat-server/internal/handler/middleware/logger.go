package middleware

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"

	chimiddleware "github.com/go-chi/chi/v5/middleware"

	log "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/logger"
)

var (
	LogEntryCtxKey = "loggerEntry"
)

func getLoggerFromRequest(req *http.Request) log.Logger {
	logger, _ := req.Context().Value(LogEntryCtxKey).(log.Logger)
	return logger
}

func wrapRequestWithLogger(req *http.Request, logger log.Logger) *http.Request {
	return req.WithContext(context.WithValue(req.Context(), LogEntryCtxKey, logger))
}

func LoggingMiddleware(logger log.Logger, level logrus.Level) Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			ww := chimiddleware.NewWrapResponseWriter(rw, req.ProtoMajor)

			t1 := time.Now()
			defer func() { // todo: enhance
				msg := fmt.Sprintf("URL: %s, Satus: %v Bytes Written: %v Header: %s Elapsed time: %s",
					req.URL, ww.Status(), ww.BytesWritten(), ww.Header(), time.Since(t1))

				logger.Logf(logrus.InfoLevel, msg)
			}()

			next.ServeHTTP(ww, wrapRequestWithLogger(req, logger))
		})
	}
}
