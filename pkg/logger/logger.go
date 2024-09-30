package logger

import (
	"context"
	"fmt"
	"net"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	Logger = zap.Logger

	Config struct {
		Mode        string
		KibanaHost  string
		KibanaPort  string
		KibanaIndex string
	}
)

func New(cfg *Config) (logger *Logger, err error) {
	switch cfg.Mode {
	case "debug":
		logger, err = zap.NewDevelopment()
		if err != nil {
			return nil, err
		}

	case "dev", "stage":
		config := zap.NewDevelopmentConfig()
		config.Level.SetLevel(zap.InfoLevel)
		config.DisableStacktrace = true

		logger, err = config.Build()
		if err != nil {
			return nil, err
		}

	case "prod":
		conn, err := net.Dial("udp", fmt.Sprintf("%s:%s", cfg.KibanaHost, cfg.KibanaPort))
		if err != nil {
			return nil, err
		}

		encodeConfig := zapcore.EncoderConfig{
			LevelKey:       "level",
			TimeKey:        "ts",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.EpochTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}

		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(encodeConfig),
			zapcore.AddSync(conn),
			zap.InfoLevel)

		logger = zap.New(
			core,
			zap.AddCaller(),
			zap.AddStacktrace(zap.FatalLevel),
		).With(zap.String("index", cfg.KibanaIndex))

	default:
		return nil, fmt.Errorf("unknown logger mode")
	}

	return logger, nil
}

func ToContext(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, Logger{}, logger)
}

func FromContext(ctx context.Context) *Logger {
	if l, ok := ctx.Value(Logger{}).(*Logger); ok {
		return l
	}
	return zap.L()
}
