package event

import (
	"golang.org/x/oauth2"
	"time"
	"fmt"
	. "Spotify_new_releases/database"
	. "Spotify_new_releases/spotify"
)

// GetFollowingArtistsAndInsertRelations get a user's following artists and insert them in the database.
func GetFollowingArtistsAndInsertRelations(dbmap *MyDbMap, userId string, token *oauth2.Token) error {
	// create new client
	client := CreateMyClientFromToken(*token)

	// get user's following artists
	artists, err := client.GetFollowingArtists(userId)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// insert the artists' data and relation between the artists and the user into the database.
	if err := dbmap.InsertArtists(artists); err != nil {
		fmt.Println(err)
		return err
	}
	timestamp := time.Now()
	if err := dbmap.InsertRelations(artists, userId, timestamp, true); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}	

// UpdateFollowingArtists is called regularly and update the relations between artists and users.
func UpdateFollowingArtists(dbmap *MyDbMap, user UserInfo) error {
	// get user info 
	userId := user.UserId

	// create new client
	client := CreateMyClientFromUserInfo(user)

	// get user's following artists
	artists, err := client.GetFollowingArtists(userId)
	if err != nil {
		fmt.Println(err)
		return err
	}

	timestamp := time.Now()

	if err := dbmap.UpdateFollowingRelation(artists, userId, timestamp); err != nil {
		fmt.Println(err)
		return err
	}

	if err := dbmap.DeleteFollowingRelations(userId, timestamp); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}