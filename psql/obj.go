package psql

import (
	"time"

	"gorm.io/gorm"
)

type (
	DataBase struct {
		Socket *gorm.DB
	}
)

type (
	ArtistEntity struct {
		gorm.Model
		Name  string
		Songs []SongEntity `gorm:"foreignKey: Artist"`
	}

	SongEntity struct {
		gorm.Model
		Artist      uint
		Name        string
		ReleaseDate time.Time
		Link        string
		Text        []CoupletEntity `gorm:"foreignKey: SongID"`
	}

	CoupletEntity struct {
		gorm.Model
		SongID     uint
		CoupletNum uint8
		Text       string
	}
)
