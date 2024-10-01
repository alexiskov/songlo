package htpsrv

import (
	"fmt"
	"net/http"
)

func New(port uint16) ServerEntity {
	return ServerEntity{Http: &http.Server{Addr: fmt.Sprintf(":%d", port)}}
}

func (srv ServerEntity) Start() (err error) {
	http.HandleFunc("/", router)
	srv.Http.ListenAndServe()
	return
}

func router(w http.ResponseWriter, r *http.Request) {

}
