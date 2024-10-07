package http_v1_handler

import (
	"go-rest-api/internal/entity"
	"go-rest-api/internal/errs"
)

type (
	Response struct {
		Description string `json:"description"`
	}

	ResponseContent struct {
		Description string         `json:"description"`
		Content     entity.Content `json:"content"`
	}
)

func Wrap(i interface{}) interface{} {
	switch v := i.(type) {
	case nil:
		return Response{Description: "ok"}

	case entity.Content:
		return ResponseContent{
			Description: "ok",
			Content:     (entity.Content)(v),
		}

	default:
		return Response{Description: v.(*errs.AppError).Msg}
	}
}
