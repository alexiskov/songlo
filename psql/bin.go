package psql

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init(host, usr, psswd, dbname string, port uint16) (err error) {
	DB, err = gorm.Open(postgres.Open(fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable ", host, usr, psswd, dbname, port)), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err == nil {
		if DB.AutoMigrate(&ArtistEntity{}, &SongEntity{}, &CoupletEntity{}) != nil {
			return
		}
	}
	return
}

//-------------------------- GETTERS -----------------------

// отдает список групп
func FindArtistByName(key string) (groupList ArtistsEnts, err error) {
	if err = DB.Where("LOWER(name) LIKE ?", "%"+key+"%").Find(&groupList).Error; err != nil {
		err = fmt.Errorf("group list finding error: %w", err)
		return
	}
	return
}

// отдает список песен группы
func (artists ArtistsEnts) GetSongs(key string, dateRelease int64, lim, off int) (count int64, response []SongEntity, err error) {
	off--
	artistIDs := make([]uint, 0, len(artists))
	for _, a := range artists {
		artistIDs = append(artistIDs, a.ID)
	}

	if dateRelease != 0 {
		if err = DB.Model(&SongEntity{}).Where("artist IN ? AND LOWER(name) LIKE ? AND release_date=?", artistIDs, "%"+key+"%", dateRelease).Count(&count).Error; err != nil {
			err = fmt.Errorf("(artist).get songs: songs count finding error: %w", err)
			return
		}
		err = DB.Where("artist IN ? AND LOWER(name) LIKE ? AND release_date=?", artistIDs, "%"+key+"%", dateRelease).Limit(lim).Offset(off).Find(&response).Error
		if err != nil {
			err = fmt.Errorf("get song from artist error: %w", err)
		}
	} else {
		if err = DB.Model(&SongEntity{}).Where("artist IN ? AND LOWER(name) LIKE ?", artistIDs, "%"+key+"%").Count(&count).Error; err != nil {
			err = fmt.Errorf("(artist).get songs: songs count without date finding error: %w", err)
			return
		}
		err = DB.Where("artist IN ? AND LOWER(name) LIKE ?", artistIDs, "%"+key+"%").Limit(lim).Offset(off).Find(&response).Error
		if err != nil {
			err = fmt.Errorf("get song from artist by release date finding error: %w", err)
		}
	}
	return
}

// отдает список песен
func FindSongs(key string, dateRelease int64, lim, off int) (count int64, response SongsEnts, err error) {
	off--
	if dateRelease != 0 {
		if err = DB.Model(&SongEntity{}).Where("LOWER(name) LIKE ? AND release_date=?", "%"+key+"%", dateRelease).Count(&count).Error; err != nil {
			err = fmt.Errorf("count of songs getting error: %w", err)
			return
		}
		if err = DB.Where("LOWER(name) LIKE ? AND release_date=?", "%"+key+"%", dateRelease).Limit(lim).Offset(off).Find(&response).Error; err != nil {
			err = fmt.Errorf("song finding by date error: %w", err)
			return
		}
	} else {
		if err = DB.Model(&SongEntity{}).Where("LOWER(name) LIKE ?", "%"+key+"%").Count(&count).Error; err != nil {
			err = fmt.Errorf("count of songs getting error: %w", err)
			return
		}
		if err = DB.Where("LOWER(name) LIKE ? ", "%"+key+"%").Limit(lim).Offset(off).Find(&response).Error; err != nil {
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

// ищет песню по тексту, опционально: []id группы,название песни, дата релиза
func FindSongByText(artistIDs []uint, key string, dateRelease int64, lim, off int) (count int64, resp SongsEnts, err error) {
	if err = DB.Model(&CoupletEntity{}).Where("LOWER(text) LIKE ?", "%"+key+"%").Count(&count).Error; err != nil {
		err = fmt.Errorf("song by text: couplet count finding error: %w", err)
		return
	}
	tempCouplets := make([]CoupletEntity, 0, count)
	if err = DB.Where("LOWER(text) LIKE ?", key).Find(&tempCouplets).Error; err != nil {
		err = fmt.Errorf("song by text: couplet finding error: %w", err)
		return
	}

	hashMapFilter := make(map[uint]bool, count)
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

	if artistIDs == nil {
		if dateRelease != 0 {
			if err = DB.Model(&SongEntity{}).Where("id IN ? AND release_date=?", songIDs, dateRelease).Count(&count).Error; err != nil {
				return 0, resp, fmt.Errorf("song by text count finding error: %w", err)
			}

			err = DB.Where("id IN ? AND release_date=?", songIDs, dateRelease).Limit(lim).Offset(off).Find(&resp).Error
			if err != nil {
				err = fmt.Errorf("song by text with date release finding error:%w", err)
			}
		} else {
			if err = DB.Model(&SongEntity{}).Where("id IN ?", songIDs).Count(&count).Error; err != nil {
				return 0, resp, fmt.Errorf("song by text count finding error: %w", err)
			}

			err = DB.Where("id IN ?", songIDs).Limit(lim).Offset(off).Find(&resp).Error
			if err != nil {
				err = fmt.Errorf("song by text without date release finding error:%w", err)
			}
		}
	} else {
		if dateRelease != 0 {
			if err = DB.Model(&SongEntity{}).Where("id IN ? AND release_date=? AND artist in ?", songIDs, dateRelease, artistIDs).Count(&count).Error; err != nil {
				return 0, resp, fmt.Errorf("song by text count finding error: %w", err)
			}

			err = DB.Where("id IN ? AND release_date=? AND artist in ?", songIDs, dateRelease, artistIDs).Limit(lim).Offset(off).Find(&resp).Error
			if err != nil {
				err = fmt.Errorf("song by text with date release finding error:%w", err)
			}
		} else {
			if err = DB.Model(&SongEntity{}).Where("id IN ? AND artist IN ?", songIDs, artistIDs).Count(&count).Error; err != nil {
				return 0, resp, fmt.Errorf("song by text count finding error: %w", err)
			}

			err = DB.Where("id IN ? AND artist IN ?", songIDs, artistIDs).Limit(lim).Offset(off).Find(&resp).Error
			if err != nil {
				err = fmt.Errorf("song by text without date release finding error:%w", err)
			}
		}
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
