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
// аргументом принимает параметры: id-исполнителя или 0, подстрока для поиска, лимит выдачи, оффсет запроса
// возвращает: общее количество элементов данных, слайс данных или ошибку
func FindArtistByName(key string, lim, off int64) (count int64, groupList []ArtistEntity, err error) {
	if err = DB.Model(&ArtistEntity{}).Where("name LIKE ?", "%"+key+"%").Count(&count).Error; err != nil {
		err = fmt.Errorf("count of artist getting error: %w", err)
		return
	}
	if err = DB.Where("name LIKE ? LIMIT ? OFFSET ?", "%"+key+"%", lim, off).Find(&groupList).Error; err != nil {
		return
	}
	return
}

// отдает список песен
// аргументом принимает: id-Исполнителя или 0, подстроку для поиска, лимит выдачи, оффсет запроса
// возвращает: общее количество элементов данных, слайс данных или ошибку
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

func (song SongEntity) GetArtist() (artist ArtistEntity, err error) {
	err = DB.Where("id=?", song.Artist).First(&artist).Error
	return
}

func (artist *ArtistEntity) GetSongs() (songs []SongEntity, err error) {
	err = DB.Where("artist=?", artist.ID).Find(&songs).Error
	return
}
