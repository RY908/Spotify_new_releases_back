package repository

import (
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/entity"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/dao"
)

func NewUserArtistsRepository() *UserArtistsRepository {
	return &UserArtistsRepository{}
}

type UserArtistsRepository struct{}

func (r *UserArtistsRepository) GetArtistsByUserID(factory dao.Factory, userID string) ([]*entity.UserArtist, error) {
	userArtistsDAO := factory.UserArtistsDAO()
	artists, err := userArtistsDAO.GetArtistsByUserID(userID)
	if err != nil {
		return nil, err
	}

	artistsEntity := entity.NewUserArtists(artists)

	return artistsEntity, nil
}
