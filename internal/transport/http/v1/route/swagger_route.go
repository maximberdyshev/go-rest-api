package http_v1_route

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SwaggerRouteRegister(ctx context.Context, r *httprouter.Router) {
	// TODO: httpSwagger.WrapHandler обернуть в middleware, для логирования
	r.HandlerFunc(http.MethodGet, "/swagger/*any", httpSwagger.WrapHandler)
}
