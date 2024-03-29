package usecase

import (
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/entity"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/service"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/spotify_service"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/dao"
	"log"
	"time"
)

func NewUpdateListeningHistoryUsecase(factory dao.Factory, logger *log.Logger, config *spotify_service.Config) *UpdateListeningHistoryUsecase {
	return &UpdateListeningHistoryUsecase{
		factory:                 factory,
		logger:                  logger,
		spotifyConfig:           config,
		artistService:           service.NewArtistService(),
		userService:             service.NewUserService(),
		listeningHistoryService: service.NewListeningHistoryService(),
	}
}

type UpdateListeningHistoryUsecase struct {
	factory                 dao.Factory
	logger                  *log.Logger
	spotifyConfig           *spotify_service.Config
	artistService           *service.ArtistService
	userService             *service.UserService
	listeningHistoryService *service.ListeningHistoryService
}

func (u *UpdateListeningHistoryUsecase) UpdateListeningHistory() error {
	u.logger.Print("Update Listening History")
	users, err := u.userService.GetAllUsers(u.factory)
	if err != nil {
		return err
	}
	for _, user := range users {
		client := spotify_service.CreateNewClientByUser(u.spotifyConfig, user)

		artists, counter, newToken, err := client.GetRecentlyPlayedArtists()
		if err != nil {
			return err
		}

		var artistsToInsert []*entity.Artist
		for _, artist := range artists {
			name, artistID, url, iconUrl := spotify_service.GetArtistInfo(artist)
			artistsToInsert = append(artistsToInsert, &entity.Artist{
				ID:      artistID,
				Name:    name,
				Url:     url,
				IconUrl: iconUrl,
			})
		}

		if err := u.artistService.InsertArtists(u.factory, artistsToInsert); err != nil {
			return err
		}

		timestamp := time.Now().UTC()
		if err := u.listeningHistoryService.InsertHistories(u.factory, artistsToInsert, user.ID, counter, false, timestamp); err != nil {
			return err
		}

		updatedUser := user.UpdateUserByToken(newToken)
		if err := u.userService.UpdateUserToken(u.factory, *updatedUser); err != nil {
			return err
		}
	}
	return nil
}
