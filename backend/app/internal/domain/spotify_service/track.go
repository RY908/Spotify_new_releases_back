package spotify_service

import (
	"fmt"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/entity"
	"github.com/zmb3/spotify"
	"strings"
	"time"
)

func (c *Client) GetNewReleases(artists []entity.Artist, userId string) ([]spotify.SimpleAlbum, error) {
	var newReleases []spotify.SimpleAlbum

	now := time.Now().UTC()
	after := now.AddDate(0, 0, -7)

	user, _ := c.client.CurrentUser()

	for _, artist := range artists {
		artistId := artist.ID
		offset := 0
		limit := 10

		opt := spotify.Options{Country: &user.Country, Limit: &limit, Offset: &offset}
		albums, err := c.client.GetArtistAlbumsOpt(spotify.ID(artistId), &opt, 2) // get albums
		if err != nil {
			err = fmt.Errorf("unable to get artist albums: %w", err)
			return nil, err
		}
		// if the album or single is released this week, add the track to newReleases
		for _, album := range albums.Albums {
			if album.ReleaseDateTime().After(after) {
				newReleases = append(newReleases, album)
			}
		}
		// time sleep is nessesary in order not to exceed spotify_service api limit
		time.Sleep(time.Millisecond * 500)

	}
	return newReleases, nil
}

// IfExclude returns if the song should be excluded from the playlist or not.
func IfExclude(user entity.User, trackName string) bool {
	res := false
	if user.IfRemixAdd == false && (strings.Contains(trackName, "Remix") || strings.Contains(trackName, "remix")) {
		res = true
	}
	if user.IfAcousticAdd == false && (strings.Contains(trackName, "Acoustic") || strings.Contains(trackName, "acoustic")) {
		res = true
	}
	return res

}
