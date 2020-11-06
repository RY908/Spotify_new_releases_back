package spotify

import (
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"time"
	"fmt"
	"Spotify_new_releases/database"
)

func (c *Client) GetCurrentUserId() (string, error) {
	user, err := c.Client.CurrentUser()
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return user.ID, nil
}

func (c *Client) CreatePlaylistForUser(userId string) (*spotify.FullPlaylist, error){
	playlist, err := c.Client.CreatePlaylistForUser(userId, "New Releases", "", true)
	if err != nil {
		fmt.Println(err)
	}
	return playlist, err
}

func (c *Client) GetRecentlyPlayedArtists() (map[spotify.ID]spotify.FullArtist, *oauth2.Token) {
	t := time.Now().UTC() // present time
	t = t.Add(-50 * time.Minute) // 50 minutes before present time 
	afterTime := t.UnixNano() / int64(time.Millisecond) // convert to milliseconds
	recentlyPlayedItems, _ := c.Client.PlayerRecentlyPlayedOpt(&spotify.RecentlyPlayedOptions{Limit: 50, AfterEpochMs: afterTime}) // get recentlyPlayedItems

	artistsSet := make(map[spotify.ID]spotify.FullArtist) // set of artists
	
	// add artist to artistsSet if not existed
	for _, item := range recentlyPlayedItems {
		for _, artist := range item.Track.Artists {
			if _, ok := artistsSet[artist.ID]; !ok {
				fullArtist, _ := c.Client.GetArtist(spotify.ID(artist.ID))
				artistsSet[artist.ID] = *fullArtist
			}
		}
	}

	token, _ := c.Client.Token()

	return artistsSet, token

}

func (c *Client) GetNewReleases(dbmap *database.MyDbMap, userId string) ([]spotify.SimpleAlbum, error) {
	var newReleases []spotify.SimpleAlbum
	// get artists from user id
	artists, err := dbmap.GetArtistsFromUserId(userId)
	if err != nil {
		fmt.Println(err)
	}

	// present UTC time
	now := time.Now().UTC()
	// 7 days ago from present time
	after := now.AddDate(0, 0, -7)
	// get current user information
	user, _ := c.Client.CurrentUser()

	for _, artist := range artists {
		artistId := artist.ArtistId // artist id to search
		offset := 0
		limit := 50 
		
		opt := spotify.Options{Country:&user.Country, Limit:&limit, Offset:&offset}
		albums, err := c.Client.GetArtistAlbumsOpt(spotify.ID(artistId), &opt, 2) // get albums
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

func ChangePlaylist(newReleases []spotify.SimpleAlbum, user database.UserInfo) error {
	// get user info 
	playlistId := user.PlaylistId
	client := CreateMyClientFromUserInfo(user).Client

	var addTracks []spotify.ID

	// ReplacePlaylistTracks
	for _, album := range newReleases {
		albumId := album.ID
		albumTracks, err := client.GetAlbumTracks(albumId)
		if err != nil {
			fmt.Println(err)
		}
		track := albumTracks.Tracks[0]

		trackId := track.ID

		addTracks = append(addTracks, trackId)
	}

	err := client.ReplacePlaylistTracks(spotify.ID(playlistId), addTracks...)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}
