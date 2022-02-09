package usecase

import (
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/entity"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/service"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/spotify_service"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/dao"
	"golang.org/x/oauth2"
)

func NewGetArtistsByUserIDUsecase(factory dao.Factory, config *spotify_service.Config) *GetArtistsByUserIDUsecase {
	return &GetArtistsByUserIDUsecase{
		factory:                 factory,
		spotifyConfig:           config,
		userService:             service.NewUserService(),
		listeningHistoryService: service.NewListeningHistoryService(),
	}
}

type GetArtistsByUserIDUsecase struct {
	factory                 dao.Factory
	spotifyConfig           *spotify_service.Config
	userService             *service.UserService
	listeningHistoryService *service.ListeningHistoryService
}

func (u *GetArtistsByUserIDUsecase) GetArtistsByUserIDUsecase(token *oauth2.Token) ([]*entity.UserArtist, error) {
	client := spotify_service.CreateNewClientByToken(u.spotifyConfig, token)

	userID, err := client.GetCurrentUserId()
	if err != nil {
		return nil, err
	}

	_, err = u.userService.GetUser(u.factory, userID)
	if err != nil {
		return nil, err
	}

	artists, err := u.listeningHistoryService.GetArtistsByUserID(u.factory, userID)
	if err != nil {
		return nil, err
	}

	return artists, nil
}
