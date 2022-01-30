package repository

import (
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/entity"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/dao"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/schema"
)

func NewArtistRepository() *ArtistRepository {
	return &ArtistRepository{}
}

type ArtistRepository struct{}

func (r *ArtistRepository) InsertArtist(factory dao.Factory, artist entity.Artist) error {
	artistDAO := factory.ArtistDAO()
	record := &schema.Artist{
		Id:      artist.ID,
		Name:    artist.Name,
		Url:     artist.Url,
		IconUrl: artist.IconUrl,
	}
	if err := artistDAO.InsertArtist(record); err != nil {
		return err
	}
	return nil

}
