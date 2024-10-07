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
		DeleteSong(string) (bool, error)
		UpdateSong(string, entity.Song) (bool, error)
		GetSongText(string, int) (entity.Content, error)
		GetFilteredSongs(entity.FilterSong, int) (entity.Content, error)
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

// GetFilteredSongs godoc
//
//	@Summary	Get filtered songs.
//	@Tags		Songs
//	@Accept		json
//	@Produce	json
//	@Param		name			query		string												false	"song name"
//	@Param		group			query		string												false	"song group"
//	@Param		release_date	query		string												false	"song release date"
//	@Param		page			query		int													false	"page"	minimum(1)
//	@Success	200				{object}	Response{content=entity.Content{items=entity.Song}}	"Success"
//	@Failure	400				{object}	Response											"Bad Request"
//	@Failure	401				{object}	Response											"Unauthorized"
//	@Failure	404				{object}	Response											"Not Found"
//	@Failure	500				{object}	Response											"Internal Server Error"
//	@Router		/songs [get]
func (h *Handler) GetFilteredSongs(w http.ResponseWriter, r *http.Request) *errs.AppError {
	page := r.URL.Query().Get("page")
	pageID, err := validatePage(page)
	if err != nil {
		h.logger.Error("Invalid page id", zap.Error(err))
		return errs.ErrBadRequest
	}

	var namePtr, groupPtr, releaseDatePTR *string
	name := r.URL.Query().Get("name")
	group := r.URL.Query().Get("group")
	releaseDate := r.URL.Query().Get("release_date")

	if name != "" {
		namePtr = &name
	}
	if group != "" {
		groupPtr = &group
	}
	if releaseDate != "" {
		releaseDatePTR = &releaseDate
	}

	filter := entity.FilterSong{
		Name:        namePtr,
		Group:       groupPtr,
		ReleaseDate: releaseDatePTR,
	}

	content, err := h.usecase.GetFilteredSongs(filter, pageID)
	if err != nil {
		h.logger.Error("Failed get filtered songs", zap.Error(err))
		if errors.Is(err, errs.ErrNotFound) {
			return errs.ErrNotFound
		}
		return errs.ErrInternal
	}

	c := entity.Content{}
	if content == c {
		h.logger.Error("Songs not found", zap.Error(err))
		return errs.ErrNotFound
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Wrap(content))
	h.logger.Info("Songs find successfully")
	return nil
}

// GetSongText godoc
//
//	@Summary	Get song text with couplet pagination.
//	@Tags		Songs
//	@Accept		json
//	@Produce	json
//	@Param		name	path		string													true	"song name"
//	@Param		page	query		int														false	"page"	minimum(1)
//	@Success	200		{object}	Response{content=entity.Content{items=entity.Couplet}}	"Success"
//	@Failure	400		{object}	Response												"Bad Request"
//	@Failure	401		{object}	Response												"Unauthorized"
//	@Failure	404		{object}	Response												"Not Found"
//	@Failure	500		{object}	Response												"Internal Server Error"
//	@Router		/songs/{name} [get]
func (h *Handler) GetSongText(w http.ResponseWriter, r *http.Request) *errs.AppError {
	params := httprouter.ParamsFromContext(r.Context())
	name := params.ByName("name")

	if err := validateName(name); err != nil {
		h.logger.Error("Validation failed", zap.Error(err))
		return errs.ErrBadRequest
	}

	page := r.URL.Query().Get("page")
	pageID, err := validatePage(page)
	if err != nil {
		h.logger.Error("Invalid page id", zap.String("song_name", name), zap.Error(err))
		return errs.ErrBadRequest
	}

	content, err := h.usecase.GetSongText(name, pageID)
	if err != nil {
		h.logger.Error("Failed get song", zap.String("song_name", name), zap.Error(err))
		if errors.Is(err, errs.ErrNotFound) {
			return errs.ErrNotFound
		}
		return errs.ErrInternal
	}

	c := entity.Content{}
	if content == c {
		h.logger.Error("Song not found", zap.String("song_name", name), zap.Error(err))
		return errs.ErrNotFound
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Wrap(content))
	h.logger.Info("Song find successfully", zap.String("song_name", name))
	return nil
}

// DeleteSong godoc
//
//	@Summary	Delete song.
//	@Tags		Songs
//	@Accept		json
//	@Produce	json
//	@Param		name	path		string		true	"song name"
//	@Success	200		{object}	Response	"Success"
//	@Failure	400		{object}	Response	"Bad Request"
//	@Failure	401		{object}	Response	"Unauthorized"
//	@Failure	404		{object}	Response	"Not Found"
//	@Failure	500		{object}	Response	"Internal Server Error"
//	@Router		/songs/{name} [delete]
func (h *Handler) DeleteSong(w http.ResponseWriter, r *http.Request) *errs.AppError {
	params := httprouter.ParamsFromContext(r.Context())
	name := params.ByName("name")

	if err := validateName(name); err != nil {
		h.logger.Error("Validation failed", zap.Error(err))
		return errs.ErrBadRequest
	}

	isDeleted, err := h.usecase.DeleteSong(name)
	if err != nil {
		h.logger.Error("Failed delete song", zap.String("song_name", name), zap.Error(err))
		if errors.Is(err, errs.ErrBadRequest) {
			return errs.ErrBadRequest
		}
		return errs.ErrInternal
	}

	if !isDeleted {
		h.logger.Error("Song not found", zap.String("song_name", name), zap.Error(err))
		return errs.ErrNotFound
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Wrap(nil))
	h.logger.Info("Song deleted successfully", zap.String("song_name", name))
	return nil
}

// UpdateSong godoc
//
//	@Summary	Update song.
//	@Tags		Songs
//	@Accept		json
//	@Produce	json
//	@Param		name	path		string		true	"song name"
//	@Param		request	body		entity.Song	true	"song text in json"
//	@Success	200		{object}	Response	"Success"
//	@Failure	400		{object}	Response	"Bad Request"
//	@Failure	401		{object}	Response	"Unauthorized"
//	@Failure	404		{object}	Response	"Not Found"
//	@Failure	500		{object}	Response	"Internal Server Error"
//	@Router		/songs/{name} [put]
func (h *Handler) UpdateSong(w http.ResponseWriter, r *http.Request) *errs.AppError {
	params := httprouter.ParamsFromContext(r.Context())
	name := params.ByName("name")

	if err := validateName(name); err != nil {
		h.logger.Error("Validation failed", zap.Error(err))
		return errs.ErrBadRequest
	}

	var updatedSong entity.Song
	if err := json.NewDecoder(r.Body).Decode(&updatedSong); err != nil {
		h.logger.Error("Invalid request payload", zap.Error(err))
		return errs.ErrBadRequest
	}
	defer r.Body.Close()

	isUpdated, err := h.usecase.UpdateSong(name, updatedSong)
	if err != nil {
		h.logger.Error("Failed update song", zap.String("song_name", name), zap.Error(err))
		if errors.Is(err, errs.ErrBadRequest) {
			return errs.ErrBadRequest
		}
		return errs.ErrInternal
	}

	if !isUpdated {
		h.logger.Error("Song not found", zap.String("song_name", name), zap.Error(err))
		return errs.ErrNotFound
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Wrap(nil))
	h.logger.Info("Song updated successfully", zap.String("song_name", name))
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
		h.logger.Error("Validation failed", zap.Error(err))
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

func validatePage(page string) (id int, err error) {
	if page != "" {
		id, err = strconv.Atoi(page)
		if err != nil {
			return 0, err
		}
		return id, nil
	}
	return 0, nil
}

func validateName(s string) error {
	if s == "" {
		return fmt.Errorf("missing or invalid query name")
	}
	return nil
}

func validateNewSong(song entity.NewSong) error {
	if song.Group == "" {
		return fmt.Errorf("missing or invalid song group")
	}
	if song.Name == "" {
		return fmt.Errorf("missing or invalid song name")
	}
	return nil
}
