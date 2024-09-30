package http_v1_handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"go-rest-api/internal/entity"
	"go-rest-api/internal/errs"
	"go-rest-api/pkg/logger"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

type (
	Usecase interface {
		AddSong(entity.NewSong) error
		DeleteSong(int) (bool, error)
		UpdateSong(int, entity.Song) (bool, error)
	}

	Handler struct {
		ctx     context.Context
		logger  *logger.Logger
		usecase Usecase
	}
)

func New(ctx context.Context, usecase Usecase) *Handler {
	return &Handler{
		ctx:     ctx,
		logger:  logger.FromContext(ctx),
		usecase: usecase,
	}
}

// TODO
// WIP..

// GetSongs godoc
//
//	@Summary	WIP..
//	@Tags		Songs
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	Response	"Success"
//	@Failure	400	{object}	Response	"Bad Request"
//	@Failure	401	{object}	Response	"Unauthorized"
//	@Failure	404	{object}	Response	"Not Found"
//	@Failure	500	{object}	Response	"Internal Server Error"
//	@Router		/songs [get]
func (h *Handler) GetSongs(w http.ResponseWriter, r *http.Request) *errs.AppError {
	return nil
}

// GetSong godoc
//
//	@Summary	WIP..
//	@Tags		Songs
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int			true	"song id"
//	@Success	200	{object}	Response	"Success"
//	@Failure	400	{object}	Response	"Bad Request"
//	@Failure	401	{object}	Response	"Unauthorized"
//	@Failure	404	{object}	Response	"Not Found"
//	@Failure	500	{object}	Response	"Internal Server Error"
//	@Router		/songs/{id} [get]
func (h *Handler) GetSong(w http.ResponseWriter, r *http.Request) *errs.AppError {
	return nil
}

// DeleteSong godoc
//
//	@Summary	Delete song.
//	@Tags		Songs
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int			true	"song id"
//	@Success	200	{object}	Response	"Success"
//	@Failure	400	{object}	Response	"Bad Request"
//	@Failure	401	{object}	Response	"Unauthorized"
//	@Failure	404	{object}	Response	"Not Found"
//	@Failure	500	{object}	Response	"Internal Server Error"
//	@Router		/songs/{id} [delete]
func (h *Handler) DeleteSong(w http.ResponseWriter, r *http.Request) *errs.AppError {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	songID, err := strconv.Atoi(id)
	if err != nil {
		h.logger.Error("Invalid song ID", zap.String("id", id), zap.Error(err))
		return errs.ErrBadRequest
	}

	if err := validateSongID(songID); err != nil {
		h.logger.Error("Song_id validation failed", zap.Error(err))
		return errs.ErrBadRequest
	}

	isDeleted, err := h.usecase.DeleteSong(songID)
	if err != nil {
		h.logger.Error("Failed delete song", zap.Int("song_id", songID), zap.Error(err))
		if errors.Is(err, errs.ErrBadRequest) {
			return errs.ErrBadRequest
		}
		return errs.ErrInternal
	}

	if !isDeleted {
		h.logger.Error("Song not found", zap.Int("song_id", songID), zap.Error(err))
		return errs.ErrNotFound
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Wrap(nil))
	h.logger.Info("Song deleted successfully", zap.Int("song_id", songID))
	return nil
}

// UpdateSong godoc
//
//	@Summary	Update song.
//	@Tags		Songs
//	@Accept		json
//	@Produce	json
//	@Param		id		path		int			true	"song id"
//	@Param		request	body		entity.Song	true	"song text in json"
//	@Success	200		{object}	Response	"Success"
//	@Failure	400		{object}	Response	"Bad Request"
//	@Failure	401		{object}	Response	"Unauthorized"
//	@Failure	404		{object}	Response	"Not Found"
//	@Failure	500		{object}	Response	"Internal Server Error"
//	@Router		/songs/{id} [put]
func (h *Handler) UpdateSong(w http.ResponseWriter, r *http.Request) *errs.AppError {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	songID, err := strconv.Atoi(id)
	if err != nil {
		h.logger.Error("Invalid song_id", zap.String("id", id), zap.Error(err))
		return errs.ErrBadRequest
	}

	if err := validateSongID(songID); err != nil {
		h.logger.Error("Song_id validation failed", zap.Error(err))
		return errs.ErrBadRequest
	}

	var updatedSong entity.Song
	if err := json.NewDecoder(r.Body).Decode(&updatedSong); err != nil {
		h.logger.Error("Invalid request payload", zap.Error(err))
		return errs.ErrBadRequest
	}
	defer r.Body.Close()

	if err := validateUpdateSong(updatedSong); err != nil {
		h.logger.Error("Song validation failed", zap.Error(err))
		return errs.ErrBadRequest
	}

	isUpdated, err := h.usecase.UpdateSong(songID, updatedSong)
	if err != nil {
		h.logger.Error("Failed update song", zap.Int("song_id", songID), zap.Error(err))
		if errors.Is(err, errs.ErrBadRequest) {
			return errs.ErrBadRequest
		}
		return errs.ErrInternal
	}

	if !isUpdated {
		h.logger.Error("Song not found", zap.Int("song_id", songID), zap.Error(err))
		return errs.ErrNotFound
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Wrap(nil))
	h.logger.Info("Song updated successfully", zap.Int("song_id", songID))
	return nil
}

// AddSong godoc
//
//	@Summary	Adding a new song.
//	@Tags		Songs
//	@Accept		json
//	@Produce	json
//	@Param		request	body		entity.NewSong	true	"json"
//	@Success	201		{object}	Response		"Success"
//	@Failure	400		{object}	Response		"Bad Request"
//	@Failure	401		{object}	Response		"Unauthorized"
//	@Failure	500		{object}	Response		"Internal Server Error"
//	@Router		/songs [post]
func (h *Handler) AddSong(w http.ResponseWriter, r *http.Request) *errs.AppError {
	var newSong entity.NewSong
	if err := json.NewDecoder(r.Body).Decode(&newSong); err != nil {
		h.logger.Error("Invalid request payload", zap.Error(err))
		return errs.ErrBadRequest
	}
	defer r.Body.Close()

	if err := validateNewSong(newSong); err != nil {
		h.logger.Error("Song validation failed", zap.Error(err))
		return errs.ErrBadRequest
	}

	if err := h.usecase.AddSong(newSong); err != nil {
		h.logger.Error("Failed to add new song", zap.Error(err))
		if errors.Is(err, errs.ErrBadRequest) {
			return errs.ErrBadRequest
		}
		return errs.ErrInternal
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Wrap(nil))
	h.logger.Info("Song added successfully")
	return nil
}

// TODO: функции ниже вероятно возможно объединить в дженерик?

func validateNewSong(song entity.NewSong) error {
	if song.Group == "" {
		return fmt.Errorf("missing or invalid song group")
	}
	if song.Name == "" {
		return fmt.Errorf("missing or invalid song name")
	}
	return nil
}

func validateSongID(id int) error {
	if id == 0 {
		return fmt.Errorf("missing or invalid song id")
	}
	return nil
}

func validateUpdateSong(song entity.Song) error {
	if song.Group == "" {
		return fmt.Errorf("missing or invalid song group")
	}
	if song.Name == "" {
		return fmt.Errorf("missing or invalid song name")
	}
	if song.ReleaseDate == "" {
		return fmt.Errorf("missing or invalid song release date")
	}
	if song.Text == "" {
		return fmt.Errorf("missing or invalid song text")
	}
	if song.Link == "" {
		return fmt.Errorf("missing or invalid song link")
	}
	return nil
}
