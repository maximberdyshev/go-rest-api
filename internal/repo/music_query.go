package repo

import (
	"fmt"
	"strings"

	"go-rest-api/internal/entity"
)

const (
	queryFindGroupID = "SELECT id FROM music_groups WHERE \"name\" = $1;"

	queryFindGroupName = "SELECT \"name\" FROM music_groups WHERE id = $1;"

	queryCreateGroup = "INSERT INTO music_groups (\"name\") VALUES ($1) RETURNING id;"

	querySaveNewSong = "INSERT INTO songs (\"name\", group_id, release_date, \"text\", \"link\") VALUES ($1, $2, $3, $4, $5);"

	queryDeleteSong = "UPDATE songs SET deleted = NOW() WHERE \"name\" = $1 AND deleted IS NULL;"

	queryGetSongText = "SELECT \"text\" FROM songs WHERE \"name\" = $1 AND deleted IS NULL;"
)

func (r *Repo) queryUpdateSong(song entity.SongDTO, name string) (string, []interface{}) {
	var str []string
	var args []interface{}
	argIndex := 1

	if song.Name != nil {
		str = append(str, fmt.Sprintf("\"name\" = $%d", argIndex))
		args = append(args, *song.Name)
		argIndex++
	}
	if song.GroupID != nil {
		str = append(str, fmt.Sprintf("group_id = $%d", argIndex))
		args = append(args, *song.GroupID)
		argIndex++
	}
	if song.ReleaseDate != nil {
		str = append(str, fmt.Sprintf("release_date = $%d", argIndex))
		args = append(args, *song.ReleaseDate)
		argIndex++
	}
	if song.Text != nil {
		str = append(str, fmt.Sprintf("\"text\" = $%d", argIndex))
		args = append(args, "{"+strings.Join(*song.Text, ",")+"}")
		argIndex++
	}
	if song.Link != nil {
		str = append(str, fmt.Sprintf("\"link\" = $%d", argIndex))
		args = append(args, *song.Link)
		argIndex++
	}

	if len(str) == 0 {
		return "", nil
	}

	where := fmt.Sprintf("\"name\" = $%d", argIndex)
	args = append(args, name)
	return "UPDATE songs SET " + strings.Join(str, ", ") + " WHERE " + where + " AND deleted IS NULL;", args
}

func (r *Repo) queryGetFilteredSongs(song entity.FilterSongDTO) (string, []interface{}) {
	baseQuery := "SELECT \"name\", group_id, release_date, \"text\", \"link\" FROM songs"
	var str []string
	var args []interface{}
	argIndex := 1

	if song.Name != nil {
		str = append(str, fmt.Sprintf("\"name\" = $%d", argIndex))
		args = append(args, *song.Name)
		argIndex++
	}
	if song.GroupID != nil {
		str = append(str, fmt.Sprintf("group_id = $%d", argIndex))
		args = append(args, *song.GroupID)
		argIndex++
	}
	if song.ReleaseDate != nil {
		str = append(str, fmt.Sprintf("release_date = $%d", argIndex))
		args = append(args, *song.ReleaseDate)
	}

	if len(str) == 0 {
		return baseQuery + " WHERE deleted IS NULL;", nil
	}

	return baseQuery + " WHERE " + strings.Join(str, " AND ") + " AND deleted IS NULL;", args
}
