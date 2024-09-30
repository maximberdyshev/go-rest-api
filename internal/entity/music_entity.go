package entity

// Models -- handlers, webapi
type (
	NewSong struct {
		Group string `json:"group" validate:"string"`
		Name  string `json:"song" validate:"string"`
	}

	SongDetail struct {
		ReleaseDate string `json:"releaseDate" validate:"string"`
		Text        string `json:"text" validate:"string"`
		Link        string `json:"link" validate:"string"`
	}

	Song struct {
		Name        string `json:"name"  validate:"string"`
		Group       string `json:"group" validate:"string"`
		ReleaseDate string `json:"release_date" validate:"string"`
		Text        string `json:"text" validate:"string"`
		Link        string `json:"link" validate:"string"`
	}
)

// DTO -- repo (postgres)
type (
	SongDTO struct {
		Name        string `json:"name"`
		GroupID     int    `json:"group_id"`
		ReleaseDate string `json:"release_date"`
		Text        string `json:"text"`
		Link        string `json:"link"`
	}
)
