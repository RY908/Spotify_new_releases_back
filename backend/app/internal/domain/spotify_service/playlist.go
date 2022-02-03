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
		return err
	}
	description := "playlist made by https://newreleases.tk"
	if err := c.client.SetPlaylistImage(playlistId, img); err != nil {
		err = fmt.Errorf("unable to set image: %w", err)
		return err
	}
	if err := c.client.ChangePlaylistDescription(playlistId, description); err != nil {
		err = fmt.Errorf("unable to change description: %w", err)
		return err
	}
	return nil

}

// ChangePlaylist change tracks in the playlist "new releases".
func (c *Client) ChangePlaylist(newReleases []spotify.SimpleAlbum, user *entity.User) error {
	playlistId := user.PlaylistID
	//client := CreateNewClientByUser(user).Client
	idSet := make(map[spotify.ID]struct{})
	pastTrackSet := make(map[spotify.ID]struct{})
	trackSet := make(map[string]struct{})
	var addTracks []spotify.ID
	var pastTracks []spotify.ID

	// get all the tracks in the playlist and put them in pastTrackSet
	playlistTrackPage, err := c.client.GetPlaylistTracks(spotify.ID(playlistId))
	if err != nil {
		err = fmt.Errorf("unable to get playlist tracks: %w", err)
		return err
	}

	// keep records of the tracks already in the playlist and delete them all
	playlistTracks := playlistTrackPage.Tracks
	for _, track := range playlistTracks {
		pastTrackSet[track.Track.ID] = struct{}{}
		pastTracks = append(pastTracks, track.Track.ID)
	}
	if _, err := c.client.RemoveTracksFromPlaylist(spotify.ID(playlistId), pastTracks...); err != nil {
		err = fmt.Errorf("unable to remove tracks: %w", err)
		return err
	}

	// retrieves track ids from newReleases. If album type is album, the first song in the album will
	// be added in the playlist.
	for _, album := range newReleases {
		albumId := album.ID
		albumTracks, err := c.client.GetAlbumTracks(albumId)
		if err != nil {
			err = fmt.Errorf("unable to get album tracks: %w", err)
			return err
		}
		fmt.Println(albumTracks.Tracks)
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
	fmt.Println(len(addTracks))

	// change the tracks in the playlist.
	change_num := (len(addTracks) - 1) / 100
	if change_num == 0 {
		if err := c.client.ReplacePlaylistTracks(spotify.ID(playlistId), addTracks...); err != nil {
			err = fmt.Errorf("unable to replace tracks in playlist: %w", err)
			return err
		}
	} else {
		if err := c.client.ReplacePlaylistTracks(spotify.ID(playlistId), addTracks[0:100]...); err != nil {
			err = fmt.Errorf("unable to replace tracks in playlist: %w", err)
			return err
		}
	}

	for i := 0; i < change_num; i++ {
		var add []spotify.ID
		fmt.Println(i)
		if i == change_num-1 {
			add = addTracks[(i+1)*100:]
		} else {
			add = addTracks[(i+1)*100 : (i+2)*100]
		}
		if _, err := c.client.AddTracksToPlaylist(spotify.ID(playlistId), add...); err != nil {
			err = fmt.Errorf("unable to add tracks to playlist: %w", err)
			return err
		}
	}

	return nil
}
