package htpsrv

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func New(port uint16) ServerEntity {
	return ServerEntity{Http: &http.Server{Addr: fmt.Sprintf(":%d", port)}}
}

func (srv ServerEntity) Start() (err error) {
	http.HandleFunc("/", router)
	err = srv.Http.ListenAndServe()
	return
}

func router(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		getProcessing(w, r)
	}
	if r.Method == http.MethodPut {

	}
	if r.Method == http.MethodPost {

	}
	if r.Method == http.MethodDelete {

	}
}

func getProcessing(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.String())
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch u.Path {
	case "/info":
		v := u.Query()

		params := URLQueryParamsEntity{Artist: v.Get("group"), Song: v.Get("song"), TextFragment: v.Get("textFragment")}

		rd := v.Get("releaseDate")
		if rd != "" {
			var t time.Time
			t, err := time.Parse("02.01.2006", rd)
			if err != nil {
				err = fmt.Errorf("query param time parsing error: %w", err)
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			params.ReleaseDate = t.Unix()
		}

		p := v.Get("page")
		if p != "" {
			pi, err := strconv.Atoi(p)
			if err != nil {
				err = fmt.Errorf("GET: query param `page` to integer type parsing error: %w", err)
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			params.Page = uint64(pi)
		}

	case "/geText":

	default:
		w.WriteHeader(http.StatusBadRequest)
	}

}
