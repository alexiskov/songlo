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
		PgCount int64              `json:"paginationCount"`
		Songs   []SongDetailEntity `json:"songs"`
	}

	SongDetailEntity struct {
		ID          uint   `yaml:"sondID"`
		Group       string `json:"group"`
		Name        string `json:"song"`
		ReleaseDate string `json:"releaseDate"`
		Link        string `json:"link"`
		Text        string `json:"text"`
	}

	SongTextEntity struct {
		SongName    string `json:"songName"`
		PgTextCount int64  `json:"coupletPaginationCount"`
		Number      int64  `json:"number"`
		Couplet     string `json:"couplet"`
	}
)
