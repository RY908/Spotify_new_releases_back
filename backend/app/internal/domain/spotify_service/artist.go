package spotify_service

import (
	"fmt"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/entity"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"time"
)

func (c *Client) GetRecentlyPlayedArtists() (map[spotify.ID]spotify.FullArtist, map[string]int, *oauth2.Token, error) {
	t := time.Now().UTC()
	t = t.Add(-20 * time.Minute)                        // 20 minutes before present time
	afterTime := t.UnixNano() / int64(time.Millisecond) // convert to milliseconds

	recentlyPlayedItems, err := c.client.PlayerRecentlyPlayedOpt(&spotify.RecentlyPlayedOptions{Limit: 50, AfterEpochMs: afterTime})
	if err != nil {
		err = fmt.Errorf("unable to get recently played tracks: %w", err)
		return nil, nil, nil, err
	}

	artistsSet := make(map[spotify.ID]spotify.FullArtist) // set of artists
	artistsCount := make(map[string]int)                  // counter

	// add an artist to artistsSet if the artist is not existed in artistsSet.
	// if already in artistsSet, increment counter.
	for _, item := range recentlyPlayedItems {
		for _, artist := range item.Track.Artists {
			if _, ok := artistsSet[artist.ID]; !ok {
				fullArtist, _ := c.client.GetArtist(spotify.ID(artist.ID))
				artistsSet[artist.ID] = *fullArtist
				artistsCount[string(artist.ID)] = 1
			} else {
				artistsCount[string(artist.ID)] += 1
			}
		}
	}

	token, _ := c.client.Token()

	return artistsSet, artistsCount, token, nil

}

// GetFollowingArtists returns artists' information a user follows.
func (c *Client) GetFollowingArtists(userId string) ([]*entity.Artist, error) {
	lastId := ""
	var artists []*entity.Artist

	// By specifying lastId, CurrentUsersFollowedArtistsOpt returns the next 50 artists the user follows.
	for {
		following, err := c.client.CurrentUsersFollowedArtistsOpt(50, lastId)
		if err != nil {
			return nil, fmt.Errorf("unable to get following artists: %w", err)
		}
		for _, following := range following.Artists {
			var name, artistId, url, iconUrl string
			name, artistId, url, iconUrl = GetArtistInfo(following)
			lastId = artistId
			artists = append(artists, &entity.Artist{ID: artistId, Name: name, Url: url, IconUrl: iconUrl})
		}

		if len(following.Artists) < 50 {
			break
		}
	}
	return artists, nil
}

// GetArtistInfo retrieves artist's name, id, url, iconUrl and returns them.
func GetArtistInfo(artist spotify.FullArtist) (string, string, string, string) {
	//var name, artistId, url, iconUrl string
	var name, artistID, url, iconUrl string
	name = artist.SimpleArtist.Name
	artistID = artist.SimpleArtist.ID.String()
	url = artist.SimpleArtist.ExternalURLs["spotify_service"]

	if len(artist.Images) > 0 {
		iconUrl = artist.Images[0].URL
	} else {
		iconUrl = ""
	}

	return name, artistID, url, iconUrl
}
