package event

import (
	"github.com/zmb3/spotify"
	//"golang.org/x/oauth2"
	"time"
	"fmt"
	. "Spotify_new_releases/database"
	. "Spotify_new_releases/spotify"
)

// This is called every 50-60 minutes to get the users recently played tracks and insert the data into the database.
func UpdateRelation(dbmap *MyDbMap) error {
	fmt.Println("Update Relation!")
	// get all the users' information from database.
	users, err := dbmap.GetAllUsers()

	if err != nil {
		err = fmt.Errorf("unable to get users: %w", err)
		return err
	}

	for _, user := range users {
		var artistsToInsert []ArtistInfo
		userId := user.UserId
		playlistId := user.PlaylistId

		client := CreateMyClientFromUserInfo(user)
		// get recently played artists.
		artists, counter, newToken, err := client.GetRecentlyPlayedArtists()
		if err != nil {
			return err
		}
		
		timestamp := time.Now()

		for _, artist := range artists {
			// get name, artistId, url, iconUrl 
			var name, artistId, url, iconUrl string
			name, artistId, url, iconUrl = GetArtistInfo(artist)
			
			artistsToInsert = append(artistsToInsert, ArtistInfo{ArtistId: artistId, Name: name, Url: url, IconUrl: iconUrl})
		}

		// update the database.
		if err := dbmap.InsertArtists(artistsToInsert); err != nil {
			err = fmt.Errorf("unable to insert artists: %w", err)
			return err
		}
		if err := dbmap.InsertRelations(artistsToInsert, counter, userId, timestamp, false); err != nil {
			err = fmt.Errorf("unable to insert artists: %w", err)
			return err
		}
		fmt.Println("userId: ", userId)
		fmt.Println("playlistId: ", playlistId)
		fmt.Println("newToken: ", newToken)
		if err := dbmap.UpdateUser(userId, playlistId, newToken); err != nil {
			err = fmt.Errorf("unable to update user: %w", err)
			return err
		}
	}
	return nil
}

// UpdatePlaylist updates the relation in the database and change the spotify playlist.
func UpdatePlaylist(dbmap *MyDbMap) error {
	fmt.Println("Update Playlist!")
	// get all the users' information from database.
	users, err := dbmap.GetAllUsers()

	if err != nil {
		err = fmt.Errorf("unable to get users: %w", err)
		return err
	}

	// for each user, get new releases and delete relations some time ago and change the spotify playlist.
	for _, user := range users {
		newReleases, err := GetNewReleasesAndDeleteRelation(dbmap, user)
		if err != nil {
			err = fmt.Errorf("unable to get new releases: %w", err)
			return err
		}
	
		if err := ChangePlaylist(newReleases, user); err != nil {
			err = fmt.Errorf("unable to change playlist: %w", err)
			return err
		}
	}
	fmt.Println("Updated")
	return nil
}

// GetNewReleasesAndDeleteRelation get artist from relation table and check if there are new tracks from the artists.
// After that, this deletes old relations from database. (currently this deletion is not available.)
func GetNewReleasesAndDeleteRelation(dbmap *MyDbMap, user UserInfo) ([]spotify.SimpleAlbum, error) {
	// get user info 
	userId := user.UserId

	// create new client
	client := CreateMyClientFromUserInfo(user)

	// get artists from user id
	artists, err := dbmap.GetArtistsFromUserId(userId)
	if err != nil {
		err = fmt.Errorf("unable to get artists from user id: %w", err)
		return nil, err
	}

	// get new releases
	newReleases, err := client.GetNewReleases(artists, userId)
	if err != nil {
		err = fmt.Errorf("unable to get new releases: %w", err)
		return nil, err
	}

	// present UTC time
	now := time.Now().UTC()
	// 7 days ago from present time
	monthAgo := now.AddDate(0, -1, 0)
	
	if err = dbmap.DeleteRelation(userId, monthAgo); err != nil {
		err = fmt.Errorf("unable to delete relation: %w", err)
		return nil, err
	}

	return newReleases, nil
}
