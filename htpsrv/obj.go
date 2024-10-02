package htpsrv

import "net/http"

type (
	ServerEntity struct {
		Http *http.Server
	}

	URLQueryParamsEntity struct {
		Group        string
		Song         string
		ReleaseDate  int64
		TextFragment string
		Page         uint64
	}

	URLQuerySongParamsEntity struct {
		SongID uint
		Page   uint
	}
)

type (
	SongRespEntity struct {
		PgCount int64            `json:"paginationCount"`
		Songs   []SongInfoEntity `json:"songs"`
	}

	SongInfoEntity struct {
		Group string `json:"group"`
		Song  string `json:"song"`
	}

	SongEntity struct {
		Group       string           `json:"group"`
		Name        string           `json:"song"`
		ReleaseDate string           `json:"releaseDate"`
		Link        string           `json:"link"`
		PgTextCount int64            `json:"coupletPaginationCount"`
		Text        []SongTextEntity `json:"text"`
	}

	SongTextEntity struct {
		Number  int64  `json:"number"`
		Couplet string `json:"couplet"`
	}
)
