package handler

import (
	"net/http"
	"strconv"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/request"
)

func GetPaginationOptsFromQuery(req *http.Request, defaultOffset int, defaultLimit int) request.PaginationOptions {
	offset, err := strconv.Atoi(req.URL.Query().Get("offset"))
	if offset == 0 || err != nil {
		offset = defaultOffset
	}

	limit, err := strconv.Atoi(req.URL.Query().Get("limit"))
	if limit == 0 || err != nil {
		limit = defaultLimit
	}

	paginationOpts := request.PaginationOptions{
		Offset: offset,
		Limit:  limit,
	}

	return paginationOpts
}
