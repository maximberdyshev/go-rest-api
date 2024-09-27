package config

import (
	"context"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		App    `yaml:"app"`
		HTTP   `yaml:"http"`
		Logger `yaml:"logger"`
		Postgres
	}

	App struct {
		Name        string `env-required:"true" yaml:"name" env:"APP_NAME"`
		Version     string `env-required:"true" yaml:"version" env:"APP_VERSION"`
		Environment string `env-required:"true" yaml:"environment" env:"ENVIRONMENT"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	Logger struct {
		Mode        string `yaml:"mode" env:"LOGGER_MODE"`
		KibanaHost  string `yaml:"kibana_host" env:"LOGGER_KIBANA_HOST"`
		KibanaPort  string `yaml:"kibana_port" env:"LOGGER_KIBANA_PORT"`
		KibanaIndex string `yaml:"kibana_index" env:"LOGGER_KIBANA_INDEX"`
	}

	Postgres struct {
		User     string `env:"POSTGRES_USER"`
		Password string `env:"POSTGRES_PASSWORD"`
		DB       string `env:"POSTGRES_DB"`
	}
)

func New() (*Config, error) {
	cfg := &Config{}
	cfgPath := "./config/config.yml"

	if err := cleanenv.ReadConfig(cfgPath, cfg); err != nil {
		return nil, err
	}

	if cfg.App.Environment == "local" {
		if err := godotenv.Load(); err != nil {
			return nil, err
		}
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func ToContext(ctx context.Context, cfg *Config) context.Context {
	return context.WithValue(ctx, Config{}, cfg)
}

func FromContext(ctx context.Context) *Config {
	return ctx.Value(Config{}).(*Config)
}
