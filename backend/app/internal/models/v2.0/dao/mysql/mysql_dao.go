package mysql

import (
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/dao"
	"github.com/go-gorp/gorp"
)

func init() {
	dao.RegisterFactory(func(db *gorp.DbMap) dao.Factory {
		return &DB{DB: db}
	})
}

type DB struct {
	DB *gorp.DbMap
}

func (d *DB) ArtistDAO() dao.Artist {
	return &artist{db: d.DB}
}

func (d *DB) ListeningHistoryDAO() dao.ListeningHistory {
	return &listeningHistory{db: d.DB}
}

func (d *DB) UserDAO() dao.User {
	return &user{db: d.DB}
}

func (d *DB) UserArtistsDAO() dao.UserArtists {
	return &userArtists{db: d.DB}
}
