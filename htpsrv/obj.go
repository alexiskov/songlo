package htpsrv

import "net/http"

type (
	ServerEntity struct {
		Http *http.Server
	}

	URLQueryParamsEntity struct {
		Artist       string
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
