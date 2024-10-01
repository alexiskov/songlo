package psql

import (
	"time"

	"gorm.io/gorm"
)

type (
	DBManager interface {
		FindArtist(adrtistID uint, key string, lim, off int64) (count int64, groupList []ArtistEntity, err error)
		FindSongs(artistID uint, key string, lim, off int64) (count int64, songList []SongEntity, err error)
	}

	DataBase struct {
		Socket *gorm.DB
	}
)

type (
	ArtistEntity struct {
		gorm.Model
		Name string
	}

	SongEntity struct {
		gorm.Model
		Artist      []ArtistEntity
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
