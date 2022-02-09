package models

import (
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/dao"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/schema"
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
	testArtist3 = &schema.Artist{
		ID:      "test_artistID3",
		Name:    "test_artist_name3",
		Url:     "test_artist_url3",
		IconUrl: "test_artist_iconUrl3",
	}
)

func InsertArtists(dao dao.Artist) error {
	artists := []*schema.Artist{
		testArtist1,
		testArtist2,
		testArtist3,
	}
	if err := dao.InsertArtists(artists); err != nil {
		return err
	}
	return nil
}
