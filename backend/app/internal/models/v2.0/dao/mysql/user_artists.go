package mysql

import (
	"fmt"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/schema"
	"github.com/go-gorp/gorp"
)

type userArtists struct {
	db *gorp.DbMap
}

func (u *userArtists) GetArtistsByUserID(userID string) ([]*schema.UserArtist, error) {
	var artists []*schema.UserArtist
	fmt.Println(userID)
	cmd := "select Artist.artistId, Artist.name, Artist.url, Artist.iconUrl, ListenTo.ifFollowing from Artist inner join ListenTo on Artist.artistId = ListenTo.artistId where ListenTo.userId = ? and ListenTo.listenCount >= 2"
	if _, err := u.db.Select(&artists, cmd, userID); err != nil {
		return nil, err
	}

	return artists, nil
}
