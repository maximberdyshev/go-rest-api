package repo

const (
	queryFindGroupID = `SELECT id FROM music_groups WHERE name = $1;`
	querySaveNewSong = `INSERT INTO songs ("name", group_id, release_date, "text", "link") VALUES ($1, $2, $3, $4, $5);`

	queryDeleteSong = `UPDATE songs SET deleted = NOW() WHERE id = $1 AND deleted IS NULL;`

	queryUpdateSong = `UPDATE songs
											SET name = $1, group_id = $2, release_date = $3, "text" = $4, "link" = $5
											WHERE id = $6
											AND deleted IS NULL;`
)
