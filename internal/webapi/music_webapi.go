package webapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"go-rest-api/config"
	"go-rest-api/internal/entity"
	"go-rest-api/internal/errs"
	"go-rest-api/pkg/logger"

	"go.uber.org/zap"
)

type Webapi struct {
	ctx    context.Context
	logger *logger.Logger
}

func New(ctx context.Context) *Webapi {
	return &Webapi{
		ctx:    ctx,
		logger: logger.FromContext(ctx),
	}
}

// GetSongDetail получает от внешнего сервиса данные о песне или возвращает ошибку.
func (wa *Webapi) GetSongDetail(newSong entity.NewSong) (songDetail entity.SongDetail, err error) {
	token := config.FromContext(wa.ctx).Webapi.Token
	webapiURL := config.FromContext(wa.ctx).Webapi.URL
	service := "info"
	externalURL := path.Join(webapiURL, service)

	params := url.Values{}
	params.Add("group", newSong.Group)
	params.Add("song", newSong.Name)

	reqURL := fmt.Sprintf("%s?%s", externalURL, params.Encode())

	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		wa.logger.Debug("Failed to create request", zap.Error(err))
		return entity.SongDetail{}, err
	}

	req.Header.Set("Authorization", token)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		wa.logger.Debug("Request to external service was executed with error", zap.String("url", req.URL.String()), zap.Error(err))
		return entity.SongDetail{}, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		if err := json.NewDecoder(res.Body).Decode(&songDetail); err != nil {
			wa.logger.Debug("Can't decode request body", zap.Error(err))
			return entity.SongDetail{}, err
		}
		return songDetail, nil
	} else {
		wa.logger.Debug("Request to external service return status code", zap.Int("status_code", res.StatusCode))
		errMsg := fmt.Sprintf("Received non-200 response status code: %s", res.Status)
		return entity.SongDetail{}, errs.NewAppError(nil, errMsg)
	}
}
