package psql

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(host, usr, psswd, dbname string, port uint16) (err error) {
	DB, err = gorm.Open(postgres.Open(fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable ", host, usr, psswd, dbname, port)), &gorm.Config{})
	if err == nil {
		if DB.AutoMigrate(&ArtistEntity{}, &SongEntity{}, &CoupletEntity{}) != nil {
			return
		}
	}
	return
}

//-------------------------- GETTERS -----------------------

// отдает список групп
func FindArtistByName(key string) (groupList []ArtistEntity, err error) {
	if err = DB.Where("name LIKE ?", "%"+key+"%").Find(&groupList).Error; err != nil {
		return
	}
	return
}

// отдает список песен группы
func (artist ArtistEntity) GetSongs(key string, dateRelease int64, lim, off int) (count int64, response []SongEntity, err error) {
	if dateRelease != 0 {
		if err = DB.Model(&SongEntity{}).Where("artist=? AND name LIKE ? AND release_date=?", artist.ID, "%"+key+"%", dateRelease).Count(&count).Error; err != nil {
			err = fmt.Errorf("(artist).get songs: songs count finding error: %w", err)
			return
		}
		err = DB.Where("artist=? AND name LIKE ? AND release_date=?", artist.ID, "%"+key+"%", dateRelease).Limit(lim).Offset(off).Find(&response).Error
	} else {
		if err = DB.Model(&SongEntity{}).Where("artist=? AND name LIKE ?", artist.ID, "%"+key+"%").Count(&count).Error; err != nil {
			err = fmt.Errorf("(artist).get songs: songs count finding error: %w", err)
			return
		}
		err = DB.Where("artist=? AND name LIKE ?", artist.ID, "%"+key+"%").Limit(lim).Offset(off).Find(&response).Error
	}
	return
}

// отдает список песен
func FindSongs(key string, dateRelease int64, lim, off int) (count int64, response SongsEnts, err error) {
	off--
	if dateRelease != 0 {
		if err = DB.Model(&SongEntity{}).Where("name LIKE ? AND release_date=?", "%"+key+"%", dateRelease).Count(&count).Error; err != nil {
			err = fmt.Errorf("count of songs getting error: %w", err)
			return
		}
		if err = DB.Where("name LIKE ? AND release_date=?", "%"+key+"%", dateRelease).Limit(lim).Offset(off).Find(&response).Error; err != nil {
			err = fmt.Errorf("song finding by date error: %w", err)
			return
		}
	} else {
		if err = DB.Model(&SongEntity{}).Where("name LIKE ?", "%"+key+"%").Count(&count).Error; err != nil {
			err = fmt.Errorf("count of songs getting error: %w", err)
			return
		}
		if err = DB.Where("name LIKE ?", "%"+key+"%").Limit(lim).Offset(off).Find(&response).Error; err != nil {
			err = fmt.Errorf("song finding without release adte error: %w", err)
			return
		}
	}

	return
}

// находит исполнителя песни
func (song SongEntity) GetArtist() (artist ArtistEntity, err error) {
	err = DB.Where("id=?", song.Artist).First(&artist).Error
	if err != nil {
		err = fmt.Errorf("artist by song getting error: %w", err)
	}
	return
}

// показывает текст песни по куплетам
func (song SongEntity) ShowText(lim, off int) (count int64, resp []CoupletEntity, err error) {
	if err = DB.Model(&CoupletEntity{}).Where("song_id=?", song.ID).Count(&count).Error; err != nil {
		err = fmt.Errorf("song couplet count finding error: %w", err)
		return
	}
	err = DB.Where("song_id=?", song.ID).Limit(lim).Offset(off).Find(&resp).Error
	if err != nil {
		err = fmt.Errorf("song text finding error: %w", err)
	}
	return
}

// ищет песню по тексту
func FindSongByText(key string, dateRelease int64, lim, off int) (resp SQlResponse, err error) {
	if err = DB.Model(&CoupletEntity{}).Where("text LIKE ?", "%"+key+"%").Count(&resp.Count).Error; err != nil {
		err = fmt.Errorf("song by text: couplet count finding error: %w", err)
		return
	}
	tempCouplets := make([]CoupletEntity, 0, resp.Count)
	if err = DB.Where("text LIKE ?", key).Find(&tempCouplets).Error; err != nil {
		err = fmt.Errorf("song by text: couplet finding error: %w", err)
		return
	}

	hashMapFilter := make(map[uint]bool, resp.Count)
	for _, c := range tempCouplets {
		if _, ok := hashMapFilter[c.SongID]; !ok {
			hashMapFilter[c.SongID] = true
		}
	}

	var songIDs []uint
	for i, v := range hashMapFilter {
		if v {
			songIDs = append(songIDs, i)
		}
	}

	if dateRelease != 0 {
		if err = DB.Model(&SongEntity{}).Where("id IN ? and release_date=?", songIDs, dateRelease).Count(&resp.Count).Error; err != nil {
			return resp, fmt.Errorf("song by text count finding error: %w", err)
		}

		resp.Data = make([]SongEntity, 0, resp.Count)

		err = DB.Where("id IN ? AND release_date=?", songIDs, dateRelease).Limit(lim).Offset(off).Find(&resp.Data).Error
	} else {
		if err = DB.Model(&SongEntity{}).Where("id IN ?", songIDs).Count(&resp.Count).Error; err != nil {
			return resp, fmt.Errorf("song by text count finding error: %w", err)
		}

		resp.Data = make([]SongEntity, 0, resp.Count)

		err = DB.Where("id IN ?", songIDs).Limit(lim).Offset(off).Find(&resp.Data).Error
	}
	return
}

//---------- SETTERS --------------

// Добавляет песню
// Агрументом принимает: Название группы, Название песни, Ссылку, Текст песни, Дату выхода
// Возвращает ошибку или nil
func AddSong(artistName, songName, link string, text string, releaseDate int64) (err error) {
	art, err := FindArtistByName(artistName)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return
		}
	}

	find := false
	for _, a := range art {
		if strings.Replace(strings.ToLower(artistName), " ", "", -1) == strings.Replace(strings.ToLower(a.Name), " ", "", -1) {
			find = true
			art = nil
			art = append(art, a)
		}
	}

	if !find {
		a := ArtistEntity{Name: artistName}
		if err = DB.Create(&a).Error; err != nil {
			return fmt.Errorf("new artist creating error: %w", err)
		}
		art = nil
		art = append(art, a)
	}

	song := SongEntity{Artist: art[0].ID, Name: songName, ReleaseDate: releaseDate, Link: link}
	if err = DB.Create(&song).Error; err != nil {
		return fmt.Errorf("new song creating error: %w", err)
	}

	if text != "" {
		tempCouplets := strings.Split(text, "\n\n")
		couplets := make([]CoupletEntity, 0, len(tempCouplets))
		for _, c := range tempCouplets {
			couplets = append(couplets, CoupletEntity{Text: c, SongID: song.ID})
		}
		if err = DB.Create(&couplets).Error; err != nil {
			return fmt.Errorf("couplets to new song adding error: %w", err)
		}
	}

	return nil
}
