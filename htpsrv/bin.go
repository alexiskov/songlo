package htpsrv

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"songlib/logger"
	"songlib/psql"
	"strconv"
	"strings"
	"time"
)

var (
	// в тз не сказано с каким шагом выполнить пагинацию и будем ли мы получать данные для нее от клиента, потому считаем что жлементов на странице 1
	SongPGstep int = 1
)

func New(port uint16) ServerEntity {
	return ServerEntity{Http: &http.Server{Addr: fmt.Sprintf(":%d", port)}}
}

func (srv ServerEntity) Start() (err error) {
	http.HandleFunc("/", router)
	logger.Log.Info(fmt.Sprintf("server  started on port  %s", srv.Http.Addr))
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
		postProcessing(w, r)
	}
	if r.Method == http.MethodDelete {

	}
}

// ------------ QUERY PROCESSING -----------

// обработка GET звпроса
func getProcessing(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.String())
	if err != nil {
		logger.Log.Debug(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	v := u.Query()

	switch u.Path {
	case "/info":
		params := URLQueryParamsEntity{Group: strings.ToLower(v.Get("group")), Song: strings.ToLower(v.Get("song")), TextFragment: strings.ToLower(v.Get("textFragment"))}

		rd := v.Get("releaseDate")
		if rd != "" {
			var t time.Time
			t, err := time.Parse("02.01.2006", rd)
			if err != nil {
				err = fmt.Errorf("query param time parsing error: %w", err)
				logger.Log.Debug(err.Error())
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
				logger.Log.Debug(err.Error())
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			params.Page = pi
		}

		if params.Page == 0 {
			params.Page = 1
		}

		sresp, err := params.SongFindingAndPrepare(SongPGstep)
		if err != nil {
			logger.Log.Debug(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		b, err := json.Marshal(sresp)
		if err != nil {
			logger.Log.Debug(err.Error())
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Write(b)

	case "/getAllTxt":

	default:
		w.WriteHeader(http.StatusBadRequest)
	}

}

// обработка PUT запроса
func putProcessing(w http.ResponseWriter, r *http.Request) {

}

// обработка POST запроса
func postProcessing(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.String())
	if err != nil {
		logger.Log.Debug(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	switch u.Path {
	case "/addsong":
		q, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Log.Debug(fmt.Errorf("addsong query reading error: %w", err).Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		song := SongDetailEntity{}
		if err = json.Unmarshal(q, &song); err != nil {
			logger.Log.Debug(fmt.Errorf("addsong query parsing error: %w", err).Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if song.Group == "" || song.Name == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		tt, err := time.Parse("02.01.2006", song.ReleaseDate)
		var rd int64 = 0
		if err != nil {
			logger.Log.Debug(fmt.Errorf("song addig: query param releaseDate parsing error: %w\n continue", err).Error())
		} else {
			rd = tt.Unix()
		}

		if err = psql.AddSong(song.Group, song.Name, song.Link, song.Text, rd); err != nil {
			logger.Log.Debug(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (queryParams URLQueryParamsEntity) SongFindingAndPrepare(limitOnPage int) (sresponse SongRespEntity, err error) {
	var c int64
	var resp psql.SongsEnts
	var groupList psql.ArtistsEnts
	if queryParams.Group == "" && queryParams.TextFragment == "" {
		c, resp, err = psql.FindSongs(queryParams.Song, queryParams.ReleaseDate, limitOnPage, limitOnPage*queryParams.Page)
		if err != nil {
			return sresponse, err
		}

	} else if queryParams.Group != "" {
		groupList, err = psql.FindArtistByName(queryParams.Group)
		if err != nil {
			return sresponse, err
		}
		c, resp, err = groupList.GetSongs(queryParams.Song, queryParams.ReleaseDate, limitOnPage, limitOnPage*queryParams.Page)
		if err != nil {
			return sresponse, err
		}
	}
	if queryParams.TextFragment != "" {
		gl := []uint{}
		for _, g := range groupList {
			gl = append(gl, g.ID)
		}
		c, resp, err = psql.FindSongByText(gl, queryParams.Song, queryParams.ReleaseDate, limitOnPage, limitOnPage*queryParams.Page)
		if err != nil {
			return sresponse, err
		}
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
	return
}
