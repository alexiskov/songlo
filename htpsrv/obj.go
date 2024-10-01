package htpsrv

import "net/http"

type (
	ServerEntity struct {
		Http *http.Server
	}
)
