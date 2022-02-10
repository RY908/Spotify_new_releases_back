package usecase

import (
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/service"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/spotify_service"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/dao"
	"golang.org/x/oauth2"
	"log"
)

func NewDeleteListeningHistoryUsecase(factory dao.Factory, logger *log.Logger, config *spotify_service.Config) *DeleteListeningHistoryUsecase {
	return &DeleteListeningHistoryUsecase{
		factory:                 factory,
		logger:                  logger,
		spotifyConfig:           config,
		userService:             service.NewUserService(),
		listeningHistoryService: service.NewListeningHistoryService(),
	}
}

type DeleteListeningHistoryUsecase struct {
	factory                 dao.Factory
	logger                  *log.Logger
	spotifyConfig           *spotify_service.Config
	userService             *service.UserService
	listeningHistoryService *service.ListeningHistoryService
}

func (u *DeleteListeningHistoryUsecase) DeleteListeningHistory(token *oauth2.Token, artistIDs []string) error {
	client := spotify_service.CreateNewClientByToken(u.spotifyConfig, token)

	userID, err := client.GetCurrentUserId()
	if err != nil {
		u.logger.Print(err)
		return err
	}

	user, err := u.userService.GetUser(u.factory, userID)
	if err != nil {
		u.logger.Print(err)
		return err
	}
	if err := u.listeningHistoryService.DeleteHistoriesByArtistIDs(u.factory, user.ID, artistIDs); err != nil {
		u.logger.Print(err)
		return err
	}
	return nil
}
