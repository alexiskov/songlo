package psql

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init(host, usr, psswd, dbname string, port int) (db DataBase, err error) {
	db.Socket, err = gorm.Open(postgres.Open(fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable ", host, usr, psswd, dbname, port)), &gorm.Config{})
	if err == nil {
		if db.Socket.AutoMigrate(&GroupEntity{}, &SongEntity{}, &CoupletEntity{}) != nil {
			return
		}
	}
	return
}

func (db *DataBase) FindAllGroups() (groupList []GroupEntity, err error) {
	err = db.Socket.Find(&groupList).Error
	return
}

func (db *DataBase) FindGroupByName(pattern string) (groupList []GroupEntity, err error) {
	err = db.Socket.Where("name LIKE ?", "%"+pattern+"%").Find(&groupList).Error
	return
}
