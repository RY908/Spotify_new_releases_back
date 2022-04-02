package usecase

import (
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/service"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/spotify_service"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/dao"
	"golang.org/x/oauth2"
	"log"
	"time"
)

func NewInitializeFollowingUsecase(factory dao.Factory, logger *log.Logger, config *spotify_service.Config) *InitializeFollowingUsecase {
	return &InitializeFollowingUsecase{
		factory:                 factory,
		logger:                  logger,
		spotifyConfig:           config,
		artistService:           service.NewArtistService(),
		userService:             service.NewUserService(),
		listeningHistoryService: service.NewListeningHistoryService(),
	}
}

type InitializeFollowingUsecase struct {
	factory                 dao.Factory
	logger                  *log.Logger
	spotifyConfig           *spotify_service.Config
	artistService           *service.ArtistService
	userService             *service.UserService
	listeningHistoryService *service.ListeningHistoryService
}

func (u *InitializeFollowingUsecase) InsertFollowingUsecase(token *oauth2.Token) error {
	client := spotify_service.CreateNewClientByToken(u.spotifyConfig, token)

	userID, err := client.GetCurrentUserId()
	if err != nil {
		u.logger.Print(err)
		return err
	}

	artists, err := client.GetFollowingArtists()
	if err != nil {
		u.logger.Print(err)
		return err
	}

	counter := map[string]int{}
	for _, artist := range artists {
		counter[artist.ID] = 1
	}

	timestamp := time.Now().UTC()
	if err := u.artistService.InsertArtists(u.factory, artists); err != nil {
		u.logger.Print(err)
		return err
	}
	if err := u.listeningHistoryService.InsertHistories(u.factory, artists, userID, counter, true, timestamp); err != nil {
		u.logger.Print(err)
		return err
	}

	return nil
}
