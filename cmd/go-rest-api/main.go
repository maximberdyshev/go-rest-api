package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"go-rest-api/config"
	"go-rest-api/internal/app"
	"go-rest-api/pkg/logger"
)

//	@title		REST-API
//	@version	1.0.0

//	@host		localhost:5000
//	@BasePath	/api/v1

//	@Security	ApiKeyAuth

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
func main() {
	srcFile := "go-rest-api/main.go"

	cfg, err := config.New()
	if err != nil {
		fmt.Printf("%s\tERROR\t%s:19\t%v\n", curTime(), srcFile, err)
		os.Exit(1)
	}
	fmt.Printf("%s\tINFO\t%s:22\tConfig loaded\n", curTime(), srcFile)

	zapLogger, err := logger.New((*logger.Config)(&cfg.Logger))
	if err != nil {
		fmt.Printf("%s\tERROR\t%s:26\t%v\n", curTime(), srcFile, err)
		os.Exit(1)
	}
	defer zapLogger.Sync()
	zapLogger.Info("Logger initialized")

	ctx := context.Background()
	ctx = config.ToContext(ctx, cfg)
	ctx = logger.ToContext(ctx, zapLogger)

	zapLogger.Info("Application launching..")
	app.Run(ctx)
}

func curTime() string {
	return time.Now().Format("2006-01-02T15:04:05.000Z0700")
}
