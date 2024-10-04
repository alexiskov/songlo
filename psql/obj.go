package psql

import (
	"gorm.io/gorm"
)

type (
	SQlResponse struct {
		Count int64
		Data  interface{}
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
		ReleaseDate int64
		Link        string
		Text        []CoupletEntity `gorm:"foreignKey: SongID"`
	}

	CoupletEntity struct {
		gorm.Model
		SongID uint
		Text   string
	}

	SongsEnts   []SongEntity
	ArtistsEnts []ArtistEntity
)
