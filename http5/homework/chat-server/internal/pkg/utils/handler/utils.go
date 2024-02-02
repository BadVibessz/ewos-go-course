package handler

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func WriteResponseAndLogError(rw http.ResponseWriter, logger *logrus.Logger, statusCode int, logMsg string, respMsg string) {
	if logMsg != "" {
		logger.Errorf(logMsg)
	}

	rw.WriteHeader(statusCode)

	if respMsg != "" {
		_, err := rw.Write([]byte(respMsg))
		if err != nil {
			logger.Errorf("error occurred writing response: %s", err)
		}
	}
}

func Paginate[T any](req *http.Request, defaultPage int, defaultLimit int, s []T) []T {
	page, _ := strconv.Atoi(req.URL.Query().Get("page"))
	if page == 0 {
		page = defaultPage
	}

	limit, _ := strconv.Atoi(req.URL.Query().Get("limit"))
	if limit == 0 {
		limit = defaultLimit
	}

	leftBound := page*limit - limit
	rightBound := leftBound + limit

	if rightBound >= len(s) {
		rightBound = len(s)
	}

	return s[leftBound:rightBound]
}
