package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"go-rest-api/internal/entity"
	"go-rest-api/internal/errs"
	"go-rest-api/pkg/logger"

	"go.uber.org/zap"
)

type Repo struct {
	ctx    context.Context
	logger *logger.Logger
	db     *sql.DB
}

func New(ctx context.Context, db *sql.DB) *Repo {
	return &Repo{
		ctx:    ctx,
		logger: logger.FromContext(ctx),
		db:     db,
	}
}

func (r *Repo) FindGroupID(group string) (id int, err error) {
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()

	err = r.db.QueryRowContext(ctx, queryFindGroupID, group).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		r.logger.Debug("Request did not return value")
		return 0, nil
	}
	if err != nil {
		r.logger.Debug("Execute sql request error", zap.Error(err))
		return 0, err
	}

	return id, nil
}

func (r *Repo) SaveNewSong(song entity.SongDTO) error {
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()

	textJSON, err := json.Marshal(song.Text)
	if err != nil {
		r.logger.Debug("Can't serialize song text to JSON", zap.Error(err))
		return err
	}

	if err = r.isDate(song.ReleaseDate); err != nil {
		r.logger.Debug("Wrong date format", zap.Error(err))
		return errs.ErrBadRequest
	}

	if _, err := r.db.ExecContext(
		ctx,
		querySaveNewSong,
		song.Name,
		song.GroupID,
		song.ReleaseDate,
		textJSON,
		song.Link,
	); err != nil {
		r.logger.Debug("Can't insert into DB", zap.Error(err))
		return err
	}

	return nil
}

func (r *Repo) DeleteSong(id int) (bool, error) {
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()

	res, err := r.db.ExecContext(ctx, queryDeleteSong, id)
	if err != nil {
		r.logger.Debug("Can't update field in table", zap.Error(err))
		return false, errs.ErrBadRequest
	}

	rows, err := res.RowsAffected()
	if err != nil {
		r.logger.Debug("Failed to get rows affected", zap.Error(err))
		return false, err
	}

	if rows == 0 {
		r.logger.Debug("Song is not exist", zap.Int("song_id", id))
	}

	return rows > 0, nil
}

func (r *Repo) UpdateSong(id int, song entity.SongDTO) (bool, error) {
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()

	textJSON, err := json.Marshal(song.Text)
	if err != nil {
		r.logger.Debug("Can't serialize song text to JSON", zap.Error(err))
		return false, err
	}

	if err = r.isDate(song.ReleaseDate); err != nil {
		r.logger.Debug("Wrong date format", zap.Error(err))
		return false, errs.ErrBadRequest
	}

	res, err := r.db.ExecContext(
		ctx,
		queryUpdateSong,
		song.Name,
		song.GroupID,
		song.ReleaseDate,
		textJSON,
		song.Link,
		id,
	)
	if err != nil {
		r.logger.Debug("Can't update field in table", zap.Error(err))
		return false, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		r.logger.Debug("Failed to get rows affected", zap.Error(err))
		return false, err
	}

	if rows == 0 {
		r.logger.Debug("Song is not exist", zap.Int("song_id", id))
	}

	return rows > 0, nil
}

func (r *Repo) isDate(str string) error {
	example := "02.01.2006"
	if _, err := time.Parse(example, str); err != nil {
		return err
	}
	return nil
}
