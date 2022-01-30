package mysql

import (
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/dao"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/schema"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/tests"
	"testing"
)

var artistInsert schema.Artist = schema.Artist{ID: "test1"}

func Test_InsertArtist(t *testing.T) {
	db, err := tests.DatabaseTestInit()
	if err != nil {
		t.Fatal(err)
	}
	dbManager := dao.NewDBManager(db)
	artistDao := dbManager.ArtistDAO()

	if err := artistDao.InsertArtist(&artistInsert); err != nil {
		t.Fatal(err)
	}
}
