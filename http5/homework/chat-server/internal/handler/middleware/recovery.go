package middleware

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"runtime/debug"
)

func RecoveryMiddleware() Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil && rvr != http.ErrAbortHandler {
					if logger := getLoggerFromRequest(req); logger != nil {
						logger.Logf(logrus.PanicLevel, "panic value: %s, stack trace: %s", rvr, debug.Stack())
					} else {
						fmt.Printf("panic occurred: %s, stack trace: %s", rvr, debug.Stack())
					}

					rw.WriteHeader(http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(rw, req)
		})
	}
}
