package usecase

import (
	"context"

	"go-rest-api/internal/entity"
	"go-rest-api/internal/errs"
	"go-rest-api/pkg/logger"

	"go.uber.org/zap"
)

type (
	Repo interface {
		FindGroupID(string) (int, error)
		CreateGroup(string) (int, error)
		CreateSong(entity.SongDTO) error
		DeleteSong(string) (bool, error)
		UpdateSong(string, entity.SongDTO) (bool, error)
		GetSongText(string) ([]string, error)
		GetFilteredSongs(entity.FilterSongDTO) ([]entity.Song, error)
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

/*
По введённым song_name и song_group:
- проверяем, что такая группа уже есть в хранилище
- если группы нет в хранилище, то она создаётся
- потом получаем данные о песне из внешнего сервиса
- записываем обогащённые данные о песне в хранилище
*/
func (uc *Usecase) AddSong(newSong entity.NewSong) error {
	groupID, err := uc.createGroup(newSong.Group)
	if err != nil {
		uc.logger.Debug("Can't create group", zap.Error(err))
		return err
	}

	songDetail, err := uc.webapi.GetSongDetail(newSong)
	if err != nil {
		uc.logger.Debug("Can't receive song detail", zap.Error(err))
		return err
	}

	songDTO := entity.SongDTO{
		Name:        &newSong.Name,
		GroupID:     &groupID,
		ReleaseDate: &songDetail.ReleaseDate,
		Text:        &songDetail.Text,
		Link:        &songDetail.Link,
	}

	if err = uc.repo.CreateSong(songDTO); err != nil {
		uc.logger.Debug("Can't save new song", zap.Error(err))
		return err
	}

	return nil
}

/*
По введённому song_name:
- "удаляем" песню из хранилища

Заметки:
1. Поиск существующей песни происходит на стороне хранилища.
2. Можно изменить удаление песни и по введённому song_id, но вшений сервис
прежде всего должен его знать.
3. Фактически, запись остаётся, но помечается отметкой об удалении.
*/
func (uc *Usecase) DeleteSong(name string) (bool, error) {
	isDeleted, err := uc.repo.DeleteSong(name)
	if err != nil {
		uc.logger.Debug("Delete song error", zap.Error(err))
		return false, err
	}

	return isDeleted, nil
}

/*
По введённому song name и new song detail:
- проверяем, что группа уже есть в хранилище
- если группы нет в хранилище, то она создаётся
- обновляем данные о песне в хранилище

Заметки:
1. Поиск существующей песни происходит на стороне хранилища.
*/
func (uc *Usecase) UpdateSong(name string, updateSong entity.Song) (bool, error) {
	groupID, err := uc.createGroup(*updateSong.Group)
	if err != nil {
		uc.logger.Debug("Can't create group", zap.Error(err))
		return false, err
	}

	song := entity.SongDTO{
		Name:        updateSong.Name,
		GroupID:     &groupID,
		ReleaseDate: updateSong.ReleaseDate,
		Text:        updateSong.Text,
		Link:        updateSong.Link,
	}

	isUpdated, err := uc.repo.UpdateSong(name, song)
	if err != nil {
		uc.logger.Debug("Update song error", zap.Error(err))
		return false, err
	}

	return isUpdated, nil
}

/*
По введённому song name и page:
- получаем текст песни
- разбиваем его на куплеты и представим, что выдаётся по 1 за раз.
Значит кол-во куплетов - это кол-во страниц, а page - указывает какую страницу
необходимо выдать на запрос.

Заметки:
1. Поиск существующей песни происходит на стороне хранилища.
*/
func (uc *Usecase) GetSongText(name string, page int) (entity.Content, error) {
	text, err := uc.repo.GetSongText(name)
	if err != nil {
		uc.logger.Debug("Find song text error", zap.Error(err))
		return entity.Content{}, err
	}

	if text == nil {
		uc.logger.Debug("Song not exist", zap.String("song_name", name))
		return entity.Content{}, errs.ErrNotFound
	}

	if page > len(text) {
		page = len(text)
	} else if page < 1 {
		page = 1
	}

	content := entity.Content{
		CurrentPage: page,
		TotalPage:   len(text),
		TotalItems:  len(text),
		Items: entity.Couplet{
			Text: text[page-1],
		},
	}

	return content, nil
}

/*
По введённым данным о песне и page:
- делаем запрос в хранилище о наличии песен с указанными параметрами
- если найдётся несколько песен, то представим, что выдаётся максимум 10 за раз.
Где каждая из вложенных структур - это структура песни.

Заметки:
1. Если будет введена группа, которой не существует, то вернётся not found.
*/
func (uc *Usecase) GetFilteredSongs(song entity.FilterSong, page int) (entity.Content, error) {
	var group string
	var groupID *int
	if song.Group != nil {
		group = *song.Group

		id, err := uc.repo.FindGroupID(group)
		if err != nil {
			uc.logger.Debug("Find group id error", zap.Error(err))
			return entity.Content{}, err
		}
		if id == 0 {
			uc.logger.Debug("Group not exist", zap.String("group", group))
			return entity.Content{}, errs.ErrNotFound
		}

		groupID = &id
	}

	s := entity.FilterSongDTO{
		Name:        song.Name,
		GroupID:     groupID,
		ReleaseDate: song.ReleaseDate,
	}

	songs, err := uc.repo.GetFilteredSongs(s)
	if err != nil {
		uc.logger.Debug("Find song error", zap.Error(err))
		return entity.Content{}, err
	}

	if songs == nil {
		uc.logger.Debug("Songs not exist")
		return entity.Content{}, errs.ErrNotFound
	}

	totalPage := len(songs) / 10
	if totalPage <= 0 {
		totalPage = 1
	}

	if page > totalPage {
		page = totalPage
	} else if page < 1 {
		page = 1
	}

	content := entity.Content{
		CurrentPage: page,
		TotalPage:   totalPage,
		TotalItems:  len(songs),
		Items:       songs,
	}

	return content, nil
}

// createGroup создаёт группу в хранилище и возвращает id записи; или возвращает ошибку.
func (uc *Usecase) createGroup(name string) (int, error) {
	groupID, err := uc.repo.FindGroupID(name)
	if err != nil {
		uc.logger.Debug("Find group id error", zap.Error(err))
		return 0, err
	}
	if groupID == 0 {
		uc.logger.Debug("Group not exist", zap.String("group", name))

		groupID, err = uc.repo.CreateGroup(name)
		if err != nil {
			uc.logger.Debug("Create group error", zap.Error(err))
			return 0, err
		}
	}
	return groupID, nil
}
