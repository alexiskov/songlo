package psql

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init(host, usr, psswd, dbname string, port int) (db DataBase, err error) {
	db.Socket, err = gorm.Open(postgres.Open(fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable ", host, usr, psswd, dbname, port)), &gorm.Config{})
	if err == nil {
		if db.Socket.AutoMigrate(&ArtistEntity{}, &SongEntity{}, &CoupletEntity{}) != nil {
			return
		}
	}
	return
}

// отдает список групп
// аргументом принимает параметры: id-исполнителя или 0, подстрока для поиска, лимит выдачи, оффсет запроса
// возвращает: общее количество элементов данных, слайс данных или ошибку
func (db *DataBase) FindArtist(adrtistID uint, key string, lim, off int64) (count int64, groupList []ArtistEntity, err error) {
	if adrtistID == 0 {
		if err = db.Socket.Model(&ArtistEntity{}).Where("name LIKE ?", "%"+key+"%").Count(&count).Error; err != nil {
			err = fmt.Errorf("count of artist getting error: %w", err)
			return
		}
		if err = db.Socket.Where("name LIKE ? LIMIT ? OFFSET ?", "%"+key+"%", lim, off).Find(&groupList).Error; err != nil {
			return
		}
	} else {
		if err = db.Socket.Model(&ArtistEntity{}).Where("id=? AND name LIKE ?", adrtistID, "%"+key+"%").Count(&count).Error; err != nil {
			err = fmt.Errorf("count of artist getting error: %w", err)
			return
		}
		if err = db.Socket.Where("id=? AND name LIKE ? LIMIT ? OFFSET ?", adrtistID, "%"+key+"%", lim, off).Find(&groupList).Error; err != nil {
			return
		}
	}
	return
}

// отдает список песен
// аргументом принимает: id-Исполнителя или 0, подстроку для поиска, лимит выдачи, оффсет запроса
// возвращает: общее количество элементов данных, слайс данных или ошибку
func (db *DataBase) FindSongs(artistID uint, key string, lim, off int64) (count int64, songList []SongEntity, err error) {
	if artistID == 0 {
		if err = db.Socket.Model(&SongEntity{}).Where("name LIKE ?", "%"+key+"%").Count(&count).Error; err != nil {
			err = fmt.Errorf("count of songs getting error: %w", err)
			return
		}
		if err = db.Socket.Where("name LIKE ? LIMIT ? OFFSET ?", "%"+key+"%", lim, off).Find(&songList).Error; err != nil {
			return
		}
	} else {
		if err = db.Socket.Model(&SongEntity{}).Where("artist=? AND name LIKE ?", artistID, "%"+key+"%").Count(&count).Error; err != nil {
			err = fmt.Errorf("count of songs getting error: %w", err)
			return
		}
		if err = db.Socket.Where("artist=? AND name LIKE ? LIMIT ? OFFSET ?", artistID, "%"+key+"%", lim, off).Find(&songList).Error; err != nil {
			return
		}
	}

	return
}
