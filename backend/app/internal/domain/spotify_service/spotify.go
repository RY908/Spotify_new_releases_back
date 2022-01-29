package spotify_service

import (
	"fmt"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/entity"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"os"
	"strings"
	"time"
)

func NewSpotifyService(config *Config) *SpotifyService {
	return &SpotifyService{
		Config: config,
	}
}

type SpotifyService struct {
	Client *spotify.Client
	Config *Config
}

func (c *Client) GetCurrentUserId() (string, error) {
	user, err := c.client.CurrentUser()
	if err != nil {
		err = fmt.Errorf("unable to get current user: %w", err)
		return "", err
	}
	return user.ID, nil
}

func (c *Client) CreatePlaylist(userId string) (*spotify.FullPlaylist, error) {
	playlist, err := c.client.CreatePlaylistForUser(userId, "New Releases", "", true)
	if err != nil {
		err = fmt.Errorf("unable to create playlist: %w", err)
		return nil, err
	}
	return playlist, err
}

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

func (c *Client) GetNewReleases(artists []entity.Artist, userId string) ([]spotify.SimpleAlbum, error) {
	var newReleases []spotify.SimpleAlbum

	now := time.Now().UTC()
	after := now.AddDate(0, 0, -7)

	user, _ := c.client.CurrentUser()

	for _, artist := range artists {
		artistId := artist.Id
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

// GetFollowingArtists returns artists' information a user follows.
func (c *Client) GetFollowingArtists(userId string) ([]entity.Artist, error) {
	lastId := ""
	var artists []entity.Artist

	// By specifying lastId, CurrentUsersFollowedArtistsOpt returns the next 50 artists the user follows.
	for {
		following, err := c.client.CurrentUsersFollowedArtistsOpt(50, lastId)
		if err != nil {
			err = fmt.Errorf("unable to get following artists: %w", err)
			return nil, err
		}
		for _, following := range following.Artists {
			var name, artistId, url, iconUrl string
			name, artistId, url, iconUrl = GetArtistInfo(following)
			lastId = artistId
			artists = append(artists, entity.Artist{Id: artistId, Name: name, Url: url, IconUrl: iconUrl})
		}

		if len(following.Artists) < 50 {
			break
		}
	}
	return artists, nil
}

// SetConfig set image and add description for a playlist.
func (c *Client) SetConfig(playlistId spotify.ID) error {
	img, _ := os.Open("./img/logo.png")
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

// GetArtistInfo retrieves artist's name, id, url, iconUrl and returns them.
func GetArtistInfo(artist spotify.FullArtist) (string, string, string, string) {
	//var name, artistId, url, iconUrl string
	var iconUrl string
	name := artist.SimpleArtist.Name
	artistId := artist.SimpleArtist.ID.String()
	url := artist.SimpleArtist.ExternalURLs["spotify_service"]
	if len(artist.Images) > 0 {
		iconUrl = artist.Images[0].URL
	} else {
		iconUrl = ""
	}

	return name, artistId, url, iconUrl
}

// ChangePlaylist change tracks in the playlist "new releases".
func (c *Client) ChangePlaylist(newReleases []spotify.SimpleAlbum, user entity.User) error {
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
