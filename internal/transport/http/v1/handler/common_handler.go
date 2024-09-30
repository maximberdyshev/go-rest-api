package http_v1_handler

import (
	"go-rest-api/internal/errs"
)

// response wrapper
type Response struct {
	Description string `json:"description"`
}

func Wrap(err *errs.AppError) Response {
	if err != nil {
		return Response{Description: err.Msg}
	}
	return Response{Description: "ok"}
}
