package psql

import (
	"fmt"

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

// отдает список групп
func FindArtistByName(key string) (groupList []ArtistEntity, err error) {
	if err = DB.Where("name LIKE ?", "%"+key+"%").Find(&groupList).Error; err != nil {
		return
	}
	return
}

// отдает список песен группы
func (artist ArtistEntity) GetSongs(key string, lim, off int64) (songs []SongEntity, err error) {
	err = DB.Where("artist=?", artist.ID).Where("name=? LIMIT ? OFFSET ?", "%"+key+"%", lim, off).Find(&songs).Error
	return
}

// отдает список песен
func FindSongsByName(key string, lim, off int64) (count int64, songList []SongEntity, err error) {
	if err = DB.Model(&SongEntity{}).Where("name LIKE ?", "%"+key+"%").Count(&count).Error; err != nil {
		err = fmt.Errorf("count of songs getting error: %w", err)
		return
	}
	if err = DB.Where("name LIKE ? LIMIT ? OFFSET ?", "%"+key+"%", lim, off).Find(&songList).Error; err != nil {
		return
	}
	return
}

// находит исполнителя песни
func (song SongEntity) GetArtist() (artist ArtistEntity, err error) {
	err = DB.Where("id=?", song.Artist).First(&artist).Error
	return
}

// показывает текст песни по куплетам
func (song SongEntity) ShowText(lim, off int64) (resp SQlResponse, err error) {
	if err = DB.Where("song_id=?", song.ID).Find(&resp.Count).Error; err != nil {
		return
	}
	resp.Data = make([]CoupletEntity, 0, lim)
	err = DB.Where("song_id=? LIMIT ? OFFSET ?", song.ID, lim, off).Find(&resp.Data).Error
	return
}

// ищет песню по тексту
func FindSongByText(key string, lim, off int64) (resp SQlResponse, err error) {
	if err = DB.Model(&CoupletEntity{}).Where("text LIKE ?", "%"+key+"%").Count(&resp.Count).Error; err != nil {
		err = fmt.Errorf("song by text couplet count finding error: %w", err)
		return
	}
	tempCouplets := make([]CoupletEntity, 0, resp.Count)
	if err = DB.Where("text LIKE ?", key).Find(&tempCouplets).Error; err != nil {
		err = fmt.Errorf("song by tex couplet finding error: %w", err)
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
	if err = DB.Model(&SongEntity{}).Where("id IN ?", songIDs).Count(&resp.Count).Error; err != nil {
		err = fmt.Errorf("song by text count finding error: %w", err)
		return
	}

	resp.Data = make([]SongEntity, 0, resp.Count)

	err = DB.Where("id IN ? LIMIT ? OFFSET ?", songIDs, lim, off).Find(&resp.Data).Error
	return
}
