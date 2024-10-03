package htpsrv

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"songlib/psql"
	"strconv"
	"time"
)

var (
	// в тз не сказано с каким шагом выполнить пагинацию и будем ли мы получать данные для нее от клиента, потому считаем что жлементов на странице 1
	SongPGstep uint64 = 1
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
		putProcessing(w, r)
	}
	if r.Method == http.MethodPost {

	}
	if r.Method == http.MethodDelete {

	}
}

// ------------ QUERY PROCESSING -----------
func getProcessing(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.String())
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	v := u.Query()

	switch u.Path {
	case "/info":
		params := URLQueryParamsEntity{Group: v.Get("group"), Song: v.Get("song"), TextFragment: v.Get("textFragment")}

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
			if pi < 1 {
				pi = 1
			}
			params.Page = uint64(pi)
		}

		sresp, err := params.SongFindingAndPrepare(SongPGstep)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		b, err := json.Marshal(sresp)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Write(b)

	case "/geText":

	default:
		w.WriteHeader(http.StatusBadRequest)
	}

}

func putProcessing(w http.ResponseWriter, r *http.Request) {

}

func (queryParams URLQueryParamsEntity) SongFindingAndPrepare(paginationDivider uint64) (sresponse SongRespEntity, err error) {
	if queryParams.Group == "" && queryParams.TextFragment == "" {
		c, resp, err := psql.FindSongs(queryParams.Song, queryParams.ReleaseDate, paginationDivider, paginationDivider*queryParams.Page)
		if err != nil {
			return sresponse, err
		}
		sresponse.PgCount = c
		for _, s := range resp {
			artist, err := s.GetArtist()
			if err != nil {
				return sresponse, err
			}

			tempDate := ""
			if s.ReleaseDate != 0 {
				tempDate = time.Unix(s.ReleaseDate, 0).Format("02.01.2006")
			} else {
				tempDate = "-"
			}

			_, c, err := s.ShowText(1, 0)
			if err != nil {
				return sresponse, err
			}
			sresponse.Songs = append(sresponse.Songs, SongDetailEntity{ID: s.ID, Group: artist.Name, Name: s.Name, Link: s.Link, ReleaseDate: tempDate, Text: c[0].Text})
		}
	}
	return
}
