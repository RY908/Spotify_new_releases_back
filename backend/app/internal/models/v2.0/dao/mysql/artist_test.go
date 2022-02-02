package mysql_test

import (
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/dao"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/schema"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/tests"
	"testing"
)

var (
	testArtist1 = &schema.Artist{
		ID:      "test_artistID1",
		Name:    "test_artist_name1",
		Url:     "test_artist_url1",
		IconUrl: "test_artist_iconUrl1",
	}
	testArtist2 = &schema.Artist{
		ID:      "test_artistID2",
		Name:    "test_artist_name2",
		Url:     "test_artist_url2",
		IconUrl: "test_artist_iconUrl2",
	}
)

func Test_InsertArtist(t *testing.T) {
	db, err := tests.DatabaseTestInit()
	if err != nil {
		t.Fatal(err)
	}
	tests.TruncateTable(t, db)
	dbManager, err := dao.NewDBManager(db)
	if err != nil {
		t.Fatal(err)
	}
	artistDao := dbManager.ArtistDAO()

	if err := artistDao.InsertArtist(testArtist1); err != nil {
		t.Fatal(err)
	}
}

func insertArtists(dao dao.Artist, artists []*schema.Artist) error {
	if err := dao.InsertArtists(artists); err != nil {
		return err
	}
	return nil
}
