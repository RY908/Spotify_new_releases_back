package main

import (
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"time"
	"fmt"
)

func getRecentlyPlayedArtists(client spotify.Client) (map[spotify.ID]spotify.FullArtist, *oauth2.Token) {
	t := time.Now().UTC() // present time
	t = t.Add(-50 * time.Minute) // 50 minutes before present time 
	afterTime := t.UnixNano() / int64(time.Millisecond) // convert to milliseconds
	recentlyPlayedItems, _ := client.PlayerRecentlyPlayedOpt(&spotify.RecentlyPlayedOptions{Limit: 50, AfterEpochMs: afterTime}) // get recentlyPlayedItems

	artistsSet := make(map[spotify.ID]spotify.FullArtist) // set of artists
	
	// add artist to artistsSet if not existed
	for _, item := range recentlyPlayedItems {
		for _, artist := range item.Track.Artists {
			if _, ok := artistsSet[artist.ID]; !ok {
				fullArtist, _ := client.GetArtist(spotify.ID(artist.ID))
				artistsSet[artist.ID] = *fullArtist
			}
		}
	}

	token, _ := client.Token()

	return artistsSet, token

}

func getNewReleases(client spotify.Client, userId string) ([]spotify.SimpleAlbum, error) {
	var newReleases []spotify.SimpleAlbum
	// get artists from user id
	artists, err := getArtistsFromUserId(userId)
	if err != nil {
		fmt.Println(err)
	}

	// present UTC time
	now := time.Now().UTC()
	// 7 days ago from present time
	after := now.AddDate(0, 0, -7)
	// get current user information
	user, _ := client.CurrentUser()

	for _, artist := range artists {
		artistId := artist.ArtistId // artist id to search
		offset := 0
		limit := 50 
		
		opt := spotify.Options{Country:&user.Country, Limit:&limit, Offset:&offset}
		albums, err := client.GetArtistAlbumsOpt(spotify.ID(artistId), &opt, 2) // get albums
		if err != nil {
			fmt.Println(err)
		}
		// if the album or single is released this week, add the track to newReleases
		for _, album := range albums.Albums {
			if album.ReleaseDateTime().After(after) {
				newReleases = append(newReleases, album)
				fmt.Println(album.Name)
			}
		}

	}
	return newReleases, nil
}