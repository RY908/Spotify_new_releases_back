package usecase

import (
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/service"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/spotify_service"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/dao"
	"golang.org/x/oauth2"
)

func NewDeleteListeningHistoryUsecase(factory dao.Factory, config *spotify_service.Config) *DeleteListeningHistoryUsecase {
	return &DeleteListeningHistoryUsecase{
		factory:                 factory,
		spotifyConfig:           config,
		userService:             service.NewUserService(),
		listeningHistoryService: service.NewListeningHistoryService(),
	}
}

type DeleteListeningHistoryUsecase struct {
	factory                 dao.Factory
	spotifyConfig           *spotify_service.Config
	userService             *service.UserService
	listeningHistoryService *service.ListeningHistoryService
}

func (u *DeleteListeningHistoryUsecase) DeleteListeningHistory(token *oauth2.Token, artistIDs []string) error {
	client := spotify_service.CreateNewClientByToken(u.spotifyConfig, token)

	userID, err := client.GetCurrentUserId()
	if err != nil {
		return err
	}

	user, err := u.userService.GetUser(u.factory, userID)
	if err != nil {
		return err
	}
	if err := u.listeningHistoryService.DeleteHistoriesByArtistIDs(u.factory, user.ID, artistIDs); err != nil {
		return err
	}
	return nil
}
