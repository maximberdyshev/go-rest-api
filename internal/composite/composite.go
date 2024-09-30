package composite

import (
	"context"
	"database/sql"

	"go-rest-api/internal/repo"
	http_v1_handler "go-rest-api/internal/transport/http/v1/handler"
	"go-rest-api/internal/usecase"
	"go-rest-api/internal/webapi"
)

type Composite struct {
	*repo.Repo
	*usecase.Usecase
	*http_v1_handler.Handler
}

func New(ctx context.Context, db *sql.DB) *Composite {
	repo := repo.New(ctx, db)
	webapi := webapi.New(ctx)
	usecase := usecase.New(ctx, repo, webapi)
	handler := http_v1_handler.New(ctx, usecase)

	return &Composite{
		Repo:    repo,
		Usecase: usecase,
		Handler: handler,
	}
}
