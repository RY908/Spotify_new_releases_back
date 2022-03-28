package models

import (
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/dao"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/schema"
	"time"
)

var (
	currentTime  = time.Now()
	testHisotry1 = &schema.ListeningHistory{
		ID:          1,
		UserID:      "test_userID1",
		ArtistID:    "test_artistID1",
		ListenCount: 1000,
		Timestamp:   currentTime,
		IsFollowing: true,
	}
	testHisotry2 = &schema.ListeningHistory{
		ID:          2,
		UserID:      "test_userID1",
		ArtistID:    "test_artistID2",
		ListenCount: 1000,
		Timestamp:   currentTime,
		IsFollowing: true,
	}
	testHisotry3 = &schema.ListeningHistory{
		ID:          3,
		UserID:      "test_userID2",
		ArtistID:    "test_artistID2",
		ListenCount: 1,
		Timestamp:   time.Date(2014, time.December, 31, 12, 13, 24, 0, time.UTC),
		IsFollowing: false,
	}
	testHisotry4 = &schema.ListeningHistory{
		ID:          4,
		UserID:      "test_userID2",
		ArtistID:    "test_artistID1",
		ListenCount: 1000,
		Timestamp:   currentTime,
		IsFollowing: false,
	}
)

func InsertHistories(dao dao.ListeningHistory) error {
	listeningHistories := []*schema.ListeningHistory{
		testHisotry1,
		testHisotry2,
		testHisotry3,
		testHisotry4,
	}
	for _, history := range listeningHistories {
		if err := dao.InsertHistory(history); err != nil {
			return err
		}
	}
	return nil
}
