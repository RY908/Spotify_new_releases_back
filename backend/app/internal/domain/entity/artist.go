package entity

import "github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/schema"

type Artist struct {
	Id      string
	Name    string
	Url     string
	IconUrl string
}

func NewArtist(artist *schema.Artist) Artist {
	return Artist{
		Id:      artist.Id,
		Name:    artist.Name,
		Url:     artist.Url,
		IconUrl: artist.IconUrl,
	}
}
