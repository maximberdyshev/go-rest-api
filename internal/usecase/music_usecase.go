package usecase

import (
	"context"
	"errors"

	"go-rest-api/internal/entity"
	"go-rest-api/internal/errs"
	"go-rest-api/pkg/logger"

	"go.uber.org/zap"
)

type (
	Repo interface {
		FindGroupID(string) (int, error)
		SaveNewSong(entity.SongDTO) error
		DeleteSong(int) (bool, error)
		UpdateSong(int, entity.SongDTO) (bool, error)
	}

	Webapi interface {
		GetSongDetail(entity.NewSong) (entity.SongDetail, error)
	}

	Usecase struct {
		ctx    context.Context
		logger *logger.Logger
		repo   Repo
		webapi Webapi
	}
)

func New(ctx context.Context, repo Repo, webapi Webapi) *Usecase {
	return &Usecase{
		ctx:    ctx,
		logger: logger.FromContext(ctx),
		repo:   repo,
		webapi: webapi,
	}
}

func (uc *Usecase) AddSong(newSong entity.NewSong) error {
	songDetail, err := uc.webapi.GetSongDetail(newSong)
	if err != nil {
		uc.logger.Debug("Can't receive song detail", zap.Error(err))
		return err
	}

	groupID, err := uc.repo.FindGroupID(newSong.Group)
	if err != nil {
		uc.logger.Debug("Find group id error", zap.Error(err))
		if errors.Is(err, errs.ErrBadRequest) {
			return errs.ErrBadRequest
		}
		return err
	}
	if groupID == 0 {
		uc.logger.Debug("Group not exist", zap.String("group", newSong.Group))
		return errs.ErrNotFound
	}

	songDTO := entity.SongDTO{
		Name:        newSong.Name,
		GroupID:     groupID,
		ReleaseDate: songDetail.ReleaseDate,
		Text:        songDetail.Text,
		Link:        songDetail.Link,
	}

	if err = uc.repo.SaveNewSong(songDTO); err != nil {
		uc.logger.Debug("Can't save new song", zap.Error(err))
		if errors.Is(err, errs.ErrBadRequest) {
			return errs.ErrBadRequest
		}
		return err
	}

	return nil
}

func (uc *Usecase) DeleteSong(id int) (bool, error) {
	isDeleted, err := uc.repo.DeleteSong(id)
	if err != nil {
		uc.logger.Debug("Delete song error", zap.Error(err))
		if errors.Is(err, errs.ErrBadRequest) {
			return false, errs.ErrBadRequest
		}
		return false, err
	}

	return isDeleted, nil
}

func (uc *Usecase) UpdateSong(id int, updateSong entity.Song) (bool, error) {
	groupID, err := uc.repo.FindGroupID(updateSong.Group)
	if err != nil {
		uc.logger.Debug("Find group id error", zap.Error(err))
		if errors.Is(err, errs.ErrBadRequest) {
			return false, errs.ErrBadRequest
		}
		return false, err
	}
	if groupID == 0 {
		uc.logger.Debug("Group not exist", zap.String("group", updateSong.Group))
		return false, nil
	}

	song := entity.SongDTO{
		Name:        updateSong.Name,
		GroupID:     groupID,
		ReleaseDate: updateSong.ReleaseDate,
		Text:        updateSong.Text,
		Link:        updateSong.Link,
	}

	isUpdated, err := uc.repo.UpdateSong(id, song)
	if err != nil {
		uc.logger.Debug("Update song error", zap.Error(err))
		if errors.Is(err, errs.ErrBadRequest) {
			return false, errs.ErrBadRequest
		}
		return false, err
	}

	return isUpdated, nil
}
