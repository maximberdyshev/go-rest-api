package http_v1_route

import (
	"context"
	"net/http"

	"go-rest-api/internal/composite"
	"go-rest-api/internal/transport/http/middleware"

	"github.com/julienschmidt/httprouter"
)

const (
	getSongs = "/api/v1/songs"
	addSong

	getSong = "/api/v1/songs/:id"
	deleteSong
	updateSong
)

func MusicRouteRegister(ctx context.Context, r *httprouter.Router, c *composite.Composite) {
	r.HandlerFunc(http.MethodGet, getSongs, middleware.Wrap(ctx, c.Handler.GetSongs))
	r.HandlerFunc(http.MethodGet, getSong, middleware.Wrap(ctx, c.Handler.GetSong))

	r.HandlerFunc(http.MethodPost, addSong, middleware.Wrap(ctx, c.Handler.AddSong))

	r.HandlerFunc(http.MethodPut, updateSong, middleware.Wrap(ctx, c.Handler.UpdateSong))

	r.HandlerFunc(http.MethodDelete, deleteSong, middleware.Wrap(ctx, c.Handler.DeleteSong))
}
