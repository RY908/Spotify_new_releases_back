package mysql_test

import (
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/dao"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/schema"
	"time"
)

var (
	testUser1 = &schema.User{
		ID:            "test_userID1",
		AccessToken:   "test_user_accessToken1",
		TokenType:     "test_user_tokenType1",
		RefreshToken:  "test_user_refreshToken1",
		Expiry:        time.Now().Add(24 * time.Hour),
		PlaylistId:    "test_user_playlistID1",
		IfRemixAdd:    true,
		IfAcousticAdd: true,
	}
	testUser2 = &schema.User{
		ID:            "test_userID2",
		AccessToken:   "test_user_accessToken2",
		TokenType:     "test_user_tokenType2",
		RefreshToken:  "test_user_refreshToken2",
		Expiry:        time.Now().Add(24 * time.Hour),
		PlaylistId:    "test_user_playlistID2",
		IfRemixAdd:    false,
		IfAcousticAdd: false,
	}
)

func insertUsers(dao dao.User, users []*schema.User) error {
	for _, user := range users {
		if err := dao.InsertUser(user); err != nil {
			return err
		}
	}
	return nil
}
