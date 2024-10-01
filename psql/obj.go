package psql

import (
	"time"

	"gorm.io/gorm"
)

type (
	DBManager interface {
		FindAllGroups() ([]GroupEntity, error)
	}

	DataBase struct {
		Socket *gorm.DB
	}
)

type (
	GroupEntity struct {
		gorm.Model
		Name string
	}

	SongEntity struct {
		gorm.Model
		Group       []GroupEntity
		Name        string
		ReleaseDate time.Time
		Link        string
	}

	CoupletEntity struct {
		gorm.Model
		SongName   []SongEntity
		CoupletNum uint8
		Text       string
	}
)
