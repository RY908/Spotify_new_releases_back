package usecase

import (
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/service"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/spotify_service"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/dao"
	"golang.org/x/oauth2"
	"log"
)

func NewUserExistsUsecase(factory dao.Factory, logger *log.Logger, config *spotify_service.Config) *UserExistsUsecase {
	return &UserExistsUsecase{
		factory:                 factory,
		logger:                  logger,
		spotifyConfig:           config,
		artistService:           service.NewArtistService(),
		userService:             service.NewUserService(),
		listeningHistoryService: service.NewListeningHistoryService(),
	}
}

type UserExistsUsecase struct {
	factory                 dao.Factory
	logger                  *log.Logger
	spotifyConfig           *spotify_service.Config
	artistService           *service.ArtistService
	userService             *service.UserService
	listeningHistoryService *service.ListeningHistoryService
}

func (u *UserExistsUsecase) IsUserExists(token *oauth2.Token) (bool, error) {
	client := spotify_service.CreateNewClientByToken(u.spotifyConfig, token)

	_, err := client.GetCurrentUserId()
	if err != nil {
		u.logger.Print(err)
		return false, err
	}
	return true, nil
}
