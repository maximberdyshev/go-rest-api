package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"go-rest-api/config"
	"go-rest-api/internal/errs"
	rw "go-rest-api/internal/transport/http/v1/handler"
	"go-rest-api/pkg/logger"

	"go.uber.org/zap"
)

type appHandler func(w http.ResponseWriter, r *http.Request) *errs.AppError

func Wrap(ctx context.Context, h appHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.FromContext(ctx)
		apiKey := config.FromContext(ctx).App.ApiKey

		w.Header().Set("Content-Type", "application/json")

		key := r.Header.Get("Authorization")
		switch key {
		case apiKey:
			logger.Info("Successful authorization", zap.String("api_key", key))

		default:
			logger.Error("Unauthorized access", zap.String("api_key", key))
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(rw.Wrap(errs.ErrUnauthorized))
			return
		}

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error("Failed to read request body", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(rw.Wrap(errs.ErrIncorrectBody))
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		logger.Info("Handling request",
			zap.String("remote_addr", r.RemoteAddr),
			zap.String("method", r.Method),
			zap.String("url", r.URL.String()),
			zap.String("body", string(bodyBytes)))

		if err := h(w, r); err != nil {
			switch err {
			// 400
			case errs.ErrBadRequest:
				logger.Error("Bad request error", zap.Error(err))
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(rw.Wrap(err))

			// 401
			case errs.ErrUnauthorized:
				logger.Error("Unauthorized error", zap.Error(err))
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(rw.Wrap(err))

			// 404
			case errs.ErrNotFound:
				logger.Error("Not found error", zap.Error(err))
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(rw.Wrap(err))

			// 500
			case errs.ErrInternal:
				logger.Error("Internal error", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(rw.Wrap(err))

			// ***
			default:
				logger.Error("Unexpected error", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(rw.Wrap(errs.ErrInternal))
			}
		}
	}
}
