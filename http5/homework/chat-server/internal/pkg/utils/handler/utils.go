package handler

import (
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/request"
	"net/http"

	headerutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/handler"
)

func GetPaginationOptsFromQuery(req *http.Request, defaultOffset int, defaultLimit int) request.PaginationOptions {
	offset, err := headerutils.GetIntParamFromQuery(req, "offset")
	if offset == 0 || err != nil {
		offset = defaultOffset
	}

	limit, err := headerutils.GetIntParamFromQuery(req, "limit")
	if limit == 0 || err != nil {
		limit = defaultLimit
	}

	paginationOpts := request.PaginationOptions{
		Offset: offset,
		Limit:  limit,
	}

	return paginationOpts
}
