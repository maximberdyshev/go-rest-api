package repo

import (
	"fmt"
	"strings"

	"go-rest-api/internal/entity"
)

const (
	queryFindGroupID = `SELECT id FROM music_groups WHERE name = $1;`

	queryFindGroupName = `SELECT "name" FROM music_groups WHERE id = $1;`

	queryCreateGroup = `INSERT INTO music_groups ("name") VALUES ($1);`

	querySaveNewSong = `INSERT INTO songs ("name", group_id, release_date, "text", "link") VALUES ($1, $2, $3, $4, $5);`

	queryDeleteSong = `UPDATE songs SET deleted = NOW() WHERE name LIKE '%$1%' AND deleted IS NULL;`

	queryUpdateSong = `UPDATE songs
											SET name = $1, group_id = $2, release_date = $3, "text" = $4, "link" = $5
											WHERE name LIKE '%$6%'
											AND deleted IS NULL;`

	queryGetSongText = `SELECT text FROM songs WHERE name LIKE '%$1%' AND deleted IS NULL;`
)

func (r *Repo) queryGetFilteredSongs(song entity.FilterSongDTO) (string, []interface{}) {
	baseQuery := "SELECT * FROM songs WHERE"
	var str []string
	var args []interface{}
	argIndex := 1

	if song.GroupID != nil {
		str = append(str, fmt.Sprintf("group_id = $%d", argIndex))
		args = append(args, *song.GroupID)
		argIndex++
	}

	if song.Name != nil {
		str = append(str, fmt.Sprintf("\"name\" LIKE '%$%d%'", argIndex))
		args = append(args, *song.Name)
		argIndex++
	}

	if song.ReleaseDate != nil {
		str = append(str, fmt.Sprintf("release_date = $%d", argIndex))
		args = append(args, *song.ReleaseDate)
	}

	if len(str) == 0 {
		return "SELECT * FROM songs", nil
	}

	return baseQuery + " " + strings.Join(str, " AND ") + " AND deleted IS NULL;", args
}
