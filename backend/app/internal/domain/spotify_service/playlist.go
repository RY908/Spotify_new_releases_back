package spotify_service

import (
	"fmt"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/entity"
	"github.com/zmb3/spotify"
	"os"
	"time"
)

func (c *Client) CreatePlaylist(userId string) (*spotify.FullPlaylist, error) {
	playlist, err := c.client.CreatePlaylistForUser(userId, "New Releases", "", true)
	if err != nil {
		err = fmt.Errorf("unable to create playlist: %w", err)
		return nil, err
	}
	return playlist, err
}

// SetConfig set image and add description for a playlist.
func (c *Client) SetConfig(playlistId spotify.ID) error {
	img, err := os.Open("./img/logo.png")
	if err != nil {
		return fmt.Errorf("unable to open image: %w", err)
	}
	description := "playlist made by https://newreleases.tk"

	if err := c.client.SetPlaylistImage(playlistId, img); err != nil {
		return fmt.Errorf("unable to set image: %w", err)
	}

	if err := c.client.ChangePlaylistDescription(playlistId, description); err != nil {
		return fmt.Errorf("unable to change description: %w", err)
	}
	return nil

}

// ChangePlaylist change tracks in the playlist "new releases".
func (c *Client) ChangePlaylist(newReleases []spotify.SimpleAlbum, user *entity.User) error {
	idSet := make(map[spotify.ID]struct{})
	pastTrackSet := make(map[spotify.ID]struct{})
	trackSet := make(map[string]struct{})

	var addTracks []spotify.ID
	var pastTracks []spotify.ID

	// get all the tracks in the playlist and put them in pastTrackSet
	playlistTrackPage, err := c.client.GetPlaylistTracks(spotify.ID(user.PlaylistID))
	if err != nil {
		return fmt.Errorf("unable to get playlist tracks: %w", err)
	}

	// keep records of the tracks already in the playlist and delete them all
	playlistTracks := playlistTrackPage.Tracks
	for _, track := range playlistTracks {
		pastTrackSet[track.Track.ID] = struct{}{}
		pastTracks = append(pastTracks, track.Track.ID)
	}
	if _, err := c.client.RemoveTracksFromPlaylist(spotify.ID(user.PlaylistID), pastTracks...); err != nil {
		return fmt.Errorf("unable to remove tracks: %w", err)
	}

	// retrieves track ids from newReleases. If album type is album, the first song in the album will
	// be added in the playlist.
	for _, album := range newReleases {
		albumTracks, err := c.client.GetAlbumTracks(album.ID)
		if err != nil {
			return fmt.Errorf("unable to get album tracks: %w", err)
		}
		track := albumTracks.Tracks[0]

		artist := track.Artists[0].Name
		trackName := track.Name
		identifier := artist + trackName // identifier is for avoiding adding both explicit song and non explicit song

		trackId := track.ID

		// avoid dupulicate tracks.
		if _, ok := idSet[trackId]; !ok {
			idSet[trackId] = struct{}{}
		} else {
			time.Sleep(time.Millisecond * 500)
			continue
		}

		// avoid adding tracks which were added last week.
		if _, ok := pastTrackSet[trackId]; ok {
			time.Sleep(time.Millisecond * 500)
			continue
		}

		// avoid adding both explicit and non explicit track
		if _, ok := trackSet[identifier]; !ok {
			trackSet[identifier] = struct{}{}
		} else {
			time.Sleep(time.Millisecond * 500)
			continue
		}

		// exclude remix and track if required
		if ok := IfExclude(user, trackName); !ok {
			addTracks = append(addTracks, trackId)
		}

		// time sleep is nessesary in order not to exceed spotify_service api limit
		time.Sleep(time.Millisecond * 500)
	}

	// change the tracks in the playlist.
	change_num := (len(addTracks) - 1) / 100
	if change_num == 0 {
		if err := c.client.ReplacePlaylistTracks(spotify.ID(user.PlaylistID), addTracks...); err != nil {
			return fmt.Errorf("unable to replace tracks in playlist: %w", err)
		}
	} else {
		if err := c.client.ReplacePlaylistTracks(spotify.ID(user.PlaylistID), addTracks[0:100]...); err != nil {
			return fmt.Errorf("unable to replace tracks in playlist: %w", err)
		}
	}

	for i := 0; i < change_num; i++ {
		var add []spotify.ID
		if i == change_num-1 {
			add = addTracks[(i+1)*100:]
		} else {
			add = addTracks[(i+1)*100 : (i+2)*100]
		}
		if _, err := c.client.AddTracksToPlaylist(spotify.ID(user.PlaylistID), add...); err != nil {
			return fmt.Errorf("unable to add tracks to playlist: %w", err)
		}
	}

	return nil
}
