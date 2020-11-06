package event

import (
	"github.com/zmb3/spotify"
	//"golang.org/x/oauth2"
	"time"
	"fmt"
	"Spotify_new_releases/database"
	"Spotify_new_releases/spotify"
)

// This is used every 50-60 minutes to get the users recently played tracks and insert the data into the database.
func updateRelation() {
	// get all users
	users, err := GetAllUsers()
	fmt.Println(users)

	if err != nil {
		fmt.Println(err)
	}

	for _, user := range users {
		// get user info 
		userId := user.UserId
		playlistId := user.PlaylistId

		client := CreateMyClientFromUserInfo(user).client
		// get recently played artists
		artists, newToken := GetRecentlyPlayedArtists(client)
		fmt.Println(newToken)
		
		timestamp := time.Now()

		for _, artist := range artists {
			// get name, artistId, url, iconUrl 
			name := artist.SimpleArtist.Name
			artistId := artist.SimpleArtist.ID.String()
			url := artist.SimpleArtist.ExternalURLs["spotify"]
			iconUrl := artist.Images[0].URL
			
			// insert artist into database
			err := InsertArtist(artistId, name, url, iconUrl)
			if err != nil {
				fmt.Println(err)
			}
	
			// insert relation into database
			err = InsertRelation(artistId, userId, timestamp)
			if err != nil {
				fmt.Println(err)
			}
	
	
		}
		
		err := UpdateUser(userId, playlistId, newToken)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func updatePlaylist() {
	// get all users
	users, err := GetAllUsers()

	if err != nil {
		fmt.Println(err)
	}

	for _, user := range users {
		newReleases := getNewReleasesAndDeleteRelation(user)
		err := ChangePlaylist(newReleases, user)

		if err != nil {
			fmt.Println(err)
		}
	}
}

func getNewReleasesAndDeleteRelation(user UserInfo) []spotify.SimpleAlbum {
	// get user info 
	userId := user.UserId

	// create new client
	client := CreateMyClientFromUserInfo(user).client

	// get new releases
	newReleases, err := GetNewReleases(client, userId)
	if err != nil {
		fmt.Println(err)
	}

	err = DeleteRelation(userId)
	if err != nil {
		fmt.Println(err)
	}

	return newReleases
}