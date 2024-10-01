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

// FindGroupID возвращает group id или ошибку.
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

// CreateGroup создаёт новую запись о группе в БД и возвращает id; или возвращает ощибку.
func (r *Repo) CreateGroup(group string) (int, error) {
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()

	res, err := r.db.ExecContext(ctx, queryCreateGroup, group)
	if err != nil {
		r.logger.Debug("Can't insert into DB", zap.Error(err))
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		r.logger.Debug("Failed to get last insert id", zap.Error(err))
		return 0, err
	}

	return int(id), nil
}

// CreateSong сохраняет песню и возвращает nil; или возвращает ошибку.
func (r *Repo) CreateSong(song entity.SongDTO) error {
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()

	textJSON, err := json.Marshal(song.Text)
	if err != nil {
		r.logger.Debug("Can't serialize song text to JSON", zap.Error(err))
		return err
	}

	if err = isDate(song.ReleaseDate); err != nil {
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

// DeleteSong удаляет песню по переданному song name и возвращает bool; или возвращает ошибку.
func (r *Repo) DeleteSong(name string) (bool, error) {
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()

	res, err := r.db.ExecContext(ctx, queryDeleteSong, name)
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
		r.logger.Debug("Song is not exist", zap.String("song_name", name))
	}

	return rows > 0, nil
}

// UpdateSong обновляет песню по переданному song name и данным и возвращает bool; или возвращает ошибку.
func (r *Repo) UpdateSong(name string, song entity.SongDTO) (bool, error) {
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()

	textJSON, err := json.Marshal(song.Text)
	if err != nil {
		r.logger.Debug("Can't serialize song text to JSON", zap.Error(err))
		return false, err
	}

	if err = isDate(song.ReleaseDate); err != nil {
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
		name,
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
		r.logger.Debug("Song is not exist", zap.String("song_name", name))
	}

	return rows > 0, nil
}

// GetSongText по song name находит песню и возвращает текст; или возвращает ошибку.
func (r *Repo) GetSongText(name string) (text []string, err error) {
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()

	var textJSON []byte
	err = r.db.QueryRowContext(ctx, queryGetSongText, name).Scan(&textJSON)
	if errors.Is(err, sql.ErrNoRows) {
		r.logger.Debug("Request did not return value")
		return nil, nil
	}
	if err != nil {
		r.logger.Debug("Execute sql request error", zap.Error(err))
		return nil, err
	}

	err = json.Unmarshal(textJSON, &text)
	if err != nil {
		r.logger.Debug("Can't deserialize JSON to song text", zap.Error(err))
		return nil, err
	}
	return text, nil
}

// GetFilteredSongs возвращает отфильтрованный список песен или ошибку.
func (r *Repo) GetFilteredSongs(song entity.FilterSongDTO) (songs []entity.Song, err error) {
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()

	if song.ReleaseDate != nil {
		if err = isDate(*song.ReleaseDate); err != nil {
			r.logger.Debug("Wrong date format", zap.Error(err))
			return nil, errs.ErrBadRequest
		}
	}

	query, args := r.queryGetFilteredSongs(song)

	rows, err := r.db.QueryContext(ctx, query, args)
	if err != nil {
		r.logger.Debug("Execute sql request error", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var song entity.SongDTO

		if err := rows.Scan(
			&song.Name,
			&song.GroupID,
			&song.ReleaseDate,
			&song.Text,
			&song.Link,
		); err != nil {
			return nil, err
		}

		group, err := r.findGroupName(song.GroupID)
		if err != nil {
			r.logger.Debug("Can't find group name", zap.Error(err))
			return nil, err
		}

		s := entity.Song{
			Name:        song.Name,
			Group:       group,
			ReleaseDate: song.ReleaseDate,
			Text:        song.Text,
			Link:        song.Link,
		}

		songs = append(songs, s)
	}

	if err = rows.Err(); err != nil {
		r.logger.Debug("Can't parse rows", zap.Error(err))
		return nil, err
	}

	return songs, nil
}

// findGroupName возвращает group name или ошибку.
func (r *Repo) findGroupName(id int) (name string, err error) {
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()

	err = r.db.QueryRowContext(ctx, queryFindGroupName, id).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		r.logger.Debug("Request did not return value")
		return "", nil
	}
	if err != nil {
		r.logger.Debug("Execute sql request error", zap.Error(err))
		return "", err
	}

	return name, nil
}

// isDate проверяет, что формат даты (DD.MM.YYYY) был указан верно.
func isDate(str string) error {
	example := "02.01.2006"
	if _, err := time.Parse(example, str); err != nil {
		return err
	}
	return nil
}
