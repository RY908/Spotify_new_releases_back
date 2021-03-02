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
		err = fmt.Errorf("unable to get following artists: %w", err)
		return err
	}

	// insert the artists' data and relation between the artists and the user into the database.
	if err := dbmap.InsertArtists(artists); err != nil {
		err = fmt.Errorf("unable to insert artists: %w", err)
		return err
	}
	timestamp := time.Now()
	if err := dbmap.InsertRelations(artists, map[string]int{}, userId, timestamp, true); err != nil {
		err = fmt.Errorf("unable to insert relations: %w", err)
		return err
	}

	return nil
}	

// UpdateFollowingArtists is called regularly and update the relations between artists and users.
func UpdateFollowingArtists(dbmap *MyDbMap) error {
	fmt.Println("update following artists")

	users, err := dbmap.GetAllUsers()

	if err != nil {
		err = fmt.Errorf("unable to get users: %w", err)
		return err
	}
	for _, user := range users {
		// get user info 
		userId := user.UserId

		// create new client
		client := CreateMyClientFromUserInfo(user)

		// get user's following artists
		artists, err := client.GetFollowingArtists(userId)
		if err != nil {
			err = fmt.Errorf("unable to get following artists: %w", err)
			return err
		}

		if err := dbmap.InsertArtists(artists); err != nil {
			err = fmt.Errorf("unable to insert artists: %w", err)
			return err
		}
		fmt.Println(artists)

		timestamp := time.Now()

		if err := dbmap.UpdateFollowingRelation(artists, userId, timestamp); err != nil {
			err = fmt.Errorf("unable to update following relation: %w", err)
			return err
		}

		if err := dbmap.DeleteFollowingRelations(userId, timestamp); err != nil {
			err = fmt.Errorf("unable to delete following relation: %w", err)
			return err
		}
	}
	fmt.Println("updated following artists")
	return nil
}
