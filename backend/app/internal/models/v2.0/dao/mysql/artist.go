package mysql

import (
	"database/sql"
	"errors"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/schema"
	"github.com/go-gorp/gorp"
)

type artist struct {
	db *gorp.DbMap
}

func (a *artist) InsertArtist(artist *schema.Artist) error {
	var insertedArtist schema.Artist
	err := a.db.SelectOne(&insertedArtist, "select * from Artist where artistId=?", artist.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if err := a.db.Insert(&schema.Artist{
				ID:      artist.ID,
				Name:    artist.Name,
				Url:     artist.Url,
				IconUrl: artist.IconUrl}); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func (a *artist) InsertArtists(artists []*schema.Artist) error {
	for _, artist := range artists {
		if err := a.InsertArtist(artist); err != nil {
			return err
		}
	}

	return nil
}
