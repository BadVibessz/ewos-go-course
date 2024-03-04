package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	myhttp "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/http"
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

func LoggingMiddleware(logger log.Logger, level log.Level) Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			ww := myhttp.NewBasicResponseWrapper(rw)

			logLev := logger.GetLevel()
			if logLev != level {
				logger.SetLevel(level)
			}

			t1 := time.Now()
			defer func() {
				// todo: colorful output?

				msg := fmt.Sprintf("URL: %s\nSatus: %v\nBytes Written: %v\nResponse: %v\nElapsed time: %s\n",
					req.URL, ww.Status(), ww.BytesWritten(), ww.Response(), time.Since(t1))

				if ww.Status() >= 400 {
					logger.Logf(log.ErrorLevel, msg)
				} else {
					logger.Logf(log.InfoLevel, msg)
				}

				// todo: needed?
				if logger.GetLevel() != logLev {
					logger.SetLevel(logLev)
				}
			}()

			next.ServeHTTP(ww, wrapRequestWithLogger(req, logger))
		})
	}
}
