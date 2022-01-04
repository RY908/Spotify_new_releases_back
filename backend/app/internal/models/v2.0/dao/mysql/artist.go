package mysql

import (
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/schema"
	"github.com/go-gorp/gorp"
)

type artist struct {
	db *gorp.DbMap
}

func (a *artist) InsertArtist(artist *schema.Artist) error {
	if err := a.db.Insert(&schema.Artist{
		Id:      artist.Id,
		Name:    artist.Name,
		Url:     artist.Url,
		IconUrl: artist.IconUrl}); err != nil {
		return err
	}
	return nil
}

func (a *artist) InsertArtists(artists []*schema.Artist) error {
	trans, err := a.db.Begin()
	if err != nil {
		return err
	}

	for _, artist := range artists {
		trans.Insert(&artist)
	}

	return trans.Commit()
}
