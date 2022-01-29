package service

import (
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/entity"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/repository"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/dao"
)

func NewArtistService() *ArtistService {
	artistRepository := repository.NewArtistRepository()
	return &ArtistService{
		repository: artistRepository,
	}
}

type ArtistService struct {
	repository *repository.ArtistRepository
}

func (s *ArtistService) InsertArtist(factory dao.Factory, artistId, name, url, iconUrl string) error {
	if err := s.repository.InsertArtist(factory, entity.Artist{
		Id:      artistId,
		Name:    name,
		Url:     url,
		IconUrl: iconUrl,
	}); err != nil {
		return err
	}
	return nil
}

func (s *ArtistService) InsertArtists(factory dao.Factory, artist []entity.Artist) error {
	for _, artist := range artist {
		if err := s.repository.InsertArtist(factory, artist); err != nil {
			return err
		}
	}
	return nil
}
