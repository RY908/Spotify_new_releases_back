package usecase

import (
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/entity"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/service"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/spotify_service"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/dao"
	"golang.org/x/oauth2"
	"log"
)

func NewInitializePlaylistUsecase(factory dao.Factory, logger *log.Logger, config *spotify_service.Config) *InitializePlaylistUsecase {
	return &InitializePlaylistUsecase{
		factory:                 factory,
		logger:                  logger,
		spotifyConfig:           config,
		userService:             service.NewUserService(),
		listeningHistoryService: service.NewListeningHistoryService(),
	}
}

type InitializePlaylistUsecase struct {
	factory                 dao.Factory
	logger                  *log.Logger
	spotifyConfig           *spotify_service.Config
	userService             *service.UserService
	listeningHistoryService *service.ListeningHistoryService
}

func (u *InitializePlaylistUsecase) CreatePlaylist(token *oauth2.Token) error {
	client := spotify_service.CreateNewClientByToken(u.spotifyConfig, token)

	userID, err := client.GetCurrentUserId()
	if err != nil {
		u.logger.Print(err)
		return err
	}

	playlist, err := client.CreatePlaylist(userID)

	if err != nil {
		u.logger.Print(err)
		return err
	}
	if err := client.SetConfig(playlist.ID); err != nil {
		u.logger.Print(err)
		return err
	}

	user := entity.NewUserCreation(userID, token, string(playlist.ID))
	if err := u.userService.InsertUser(u.factory, *user); err != nil {
		u.logger.Print(err)
		return err
	}

	return nil
}
