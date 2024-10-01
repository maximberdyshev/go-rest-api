package entity

// Models -- handlers, webapi
type (
	// add new song
	NewSong struct {
		Group string `json:"group" validate:"string"`
		Name  string `json:"song" validate:"string"`
	}

	// get song detail from external api
	SongDetail struct {
		ReleaseDate string   `json:"releaseDate" validate:"string"`
		Text        []string `json:"text" validate:"array"`
		Link        string   `json:"link" validate:"string"`
	}

	// update, filtered song
	Song struct {
		Name        string   `json:"name"  validate:"string"`
		Group       string   `json:"group" validate:"string"`
		ReleaseDate string   `json:"release_date" validate:"string"`
		Text        []string `json:"text" validate:"array"`
		Link        string   `json:"link" validate:"string"`
	}

	// filtered songs
	FilterSong struct {
		Name        *string `json:"name,omitempty"  validate:"string"`
		Group       *string `json:"group,omitempty" validate:"string"`
		ReleaseDate *string `json:"release_date,omitempty" validate:"string"`
	}
)

// Models -- response
type (
	Content struct {
		CurrentPage int         `json:"current_page"`
		TotalPage   int         `json:"total_page"`
		TotalItems  int         `json:"total_items"`
		Items       interface{} `json:"items"`
	}

	Couplet struct {
		Text string `json:"text"`
	}
)

// DTO -- repo (postgres)
type (
	SongDTO struct {
		Name        string
		GroupID     int
		ReleaseDate string
		Text        []string
		Link        string
	}

	FilterSongDTO struct {
		Name        *string
		GroupID     *int
		ReleaseDate *string
	}
)
