package entity

import "github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/schema"

type Artist struct {
	ID      string
	Name    string
	Url     string
	IconUrl string
}

func NewArtist(artist *schema.Artist) Artist {
	return Artist{
		ID:      artist.Id,
		Name:    artist.Name,
		Url:     artist.Url,
		IconUrl: artist.IconUrl,
	}
}

func NewArtists(artists []schema.Artist) []Artist {
	var artistsEntity []Artist
	for _, artist := range artists {
		artistsEntity = append(artistsEntity, NewArtist(&artist))
	}
	return artistsEntity
}
