package postgres

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

const (
	_defaultMaxLifetime     = 60 * time.Second
	_defaultConnMaxIdleTime = 30 * time.Second
	_defaultMaxIdleConns    = 5
	_defaultMaxOpenConns    = 10
)

type Config struct {
	Host     string
	Port     string
	DB       string
	User     string
	Password string
}

func New(cfg *Config) (*sql.DB, error) {
	src := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DB,
	)

	db, err := sql.Open("postgres", src)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(_defaultMaxLifetime)
	db.SetConnMaxIdleTime(_defaultConnMaxIdleTime)
	db.SetMaxIdleConns(_defaultMaxIdleConns)
	db.SetMaxOpenConns(_defaultMaxOpenConns)

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, err
}
