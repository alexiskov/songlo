package htpsrv

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"songlib/logger"
	"songlib/psql"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"
)

var (
	// в тз не сказано с каким шагом выполнить пагинацию и будем ли мы получать данные для нее от клиента, потому считаем что жлементов на странице 1
	SongPGstep     int = 1
	SongTextPGstep int = 1
)

func Start(port string) (err error) {
	r := chi.NewRouter()
	r.Get("/info", getProcessing)
	r.Get("/getAllTxt", bezModaKorocheBylo)
	r.Post("/addsong", songAdd)
	r.Put("/updateSong", putProcessing)
	r.Delete("/delete", deleteSong)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
	))
	logger.Log.Info(fmt.Sprint("server  started"))
	err = http.ListenAndServe(":"+port, r)
	return
}

func deleteSong(w http.ResponseWriter, r *http.Request) {
	if j, err := io.ReadAll(r.Body); err != nil {
		logger.Log.Debug(fmt.Errorf("query on deleting song reading error: %w", err).Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		delSongQr := URLQuerySongParamsEntity{}
		if err = json.Unmarshal(j, &delSongQr); err != nil {
			logger.Log.Debug(fmt.Errorf("query on deleting song json parsing error: %w", err).Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err = psql.Remove(delSongQr.SongID); err != nil {
			logger.Log.Debug(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			w.WriteHeader(http.StatusOK)
		}

	}
}

// ------------ QUERY PROCESSING -----------

// обработка GET звпроса
func getProcessing(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()

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
}

func bezModaKorocheBylo(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	songID, err := strconv.Atoi(v.Get("id"))
	if err != nil {
		logger.Log.Debug(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	pageNum, err := strconv.Atoi(v.Get("page"))
	if err != nil {
		logger.Log.Debug(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	queryParams := URLQuerySongParamsEntity{SongID: uint(songID), Page: uint(pageNum)}
	if queryParams.SongID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		logger.Log.Debug("get all text input query param songID == 0, expected >0")
		return
	}

	t, err := queryParams.SongTextFindingAndPrepare()
	if err != nil {
		logger.Log.Debug(err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if b, err := json.Marshal(t); err != nil {
		logger.Log.Debug(fmt.Errorf("song text json marshaling error: %w", err).Error())
		w.WriteHeader(http.StatusNotFound)
		return
	} else {
		w.Write(b)
	}
}

// обработка PUT запроса
func putProcessing(w http.ResponseWriter, r *http.Request) {
	readByte, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Log.Debug(fmt.Errorf("updateSong body of query reading error: %w", err).Error())
		return
	}

	song := SongDetailEntity{}
	if err = json.Unmarshal(readByte, &song); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Log.Debug(fmt.Errorf("update song body of qery parsing error: %w", err).Error())
		return
	}

	if err = song.UpdateSong(); err != nil {
		logger.Log.Debug(err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

// обработка POST запроса
func songAdd(w http.ResponseWriter, r *http.Request) {
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

}

// находит песню и возвращает данные
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

// берет текст для песни по id из базы
func (queryParams URLQuerySongParamsEntity) SongTextFindingAndPrepare() (couplets SongTextEntity, err error) {
	song := psql.SongEntity{}
	song.ID = queryParams.SongID
	c, resp, err := song.ShowText(SongTextPGstep, SongTextPGstep*int(queryParams.Page))
	if err != nil {
		return
	}
	couplets.PgCount = c
	for _, t := range resp {
		couplets.Couplet += "\n\n" + t.Text
	}
	return
}

// обновляет данные песни
func (queryParams SongDetailEntity) UpdateSong() (err error) {
	var tt time.Time
	if queryParams.ReleaseDate != "" {
		tt, err = time.Parse("02.01.2006", queryParams.ReleaseDate)
		if err != nil {
			return fmt.Errorf("update song: time parsing error: %w", err)
		}
	}
	song := psql.SongEntity{Name: queryParams.Name, ReleaseDate: tt.Unix(), Link: queryParams.Link}
	tsTextArr := strings.Split(queryParams.Text, "\n\n")
	for _, txt := range tsTextArr {
		song.Text = append(song.Text, psql.CoupletEntity{Text: txt})
	}

	err = song.Update(queryParams.Group, queryParams.ID)
	return
}
