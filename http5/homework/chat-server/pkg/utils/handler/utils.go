package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

func WriteErrResponseAndLog(rw http.ResponseWriter, logger *logrus.Logger, statusCode int, logMsg string, respMsg string) {
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

func GetOffsetAndLimitFromQuery(req *http.Request, defaultOffset int, defaultLimit int) (off int, lim int) {
	offset, err := strconv.Atoi(req.URL.Query().Get("offset"))
	if offset == 0 || err != nil {
		offset = defaultOffset
	}

	limit, err := strconv.Atoi(req.URL.Query().Get("limit"))
	if limit == 0 || err != nil {
		limit = defaultLimit
	}

	return offset, limit
}

var (
	ErrNoHeaderProvided      = errors.New("no header provided")
	ErrInvalidHeaderProvided = errors.New("invalid header provided")
)

func GetIntHeaderByKey(req *http.Request, key string) (int, error) {
	str := req.Header.Get(key)
	if str == "" {
		return -1, ErrNoHeaderProvided
	}

	val, err := strconv.Atoi(str)
	if err != nil {
		return -1, ErrInvalidHeaderProvided
	}

	return val, nil
}
