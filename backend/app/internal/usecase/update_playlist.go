package usecase

import (
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/service"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/spotify_service"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/dao"
	"time"
)

func NewUpdatePlaylistUsecase(factory dao.Factory, config *spotify_service.Config) *UpdatePlaylistUsecase {
	return &UpdatePlaylistUsecase{
		factory:                 factory,
		spotifyConfig:           config,
		artistService:           service.NewArtistService(),
		userService:             service.NewUserService(),
		listeningHistoryService: service.NewListeningHistoryService(),
	}
}

type UpdatePlaylistUsecase struct {
	factory                 dao.Factory
	spotifyConfig           *spotify_service.Config
	artistService           *service.ArtistService
	userService             *service.UserService
	listeningHistoryService *service.ListeningHistoryService
}

func (u *UpdatePlaylistUsecase) UpdatePlaylistHistory() error {
	users, err := u.userService.GetAllUsers(u.factory)
	if err != nil {
		return err
	}
	for _, user := range users {
		client := spotify_service.CreateNewClientByUser(u.spotifyConfig, user)

		artists, err := u.listeningHistoryService.GetArtistsByUserID(u.factory, user.ID)
		if err != nil {
			return err
		}

		newReleases, err := client.GetNewReleases(artists, user.ID)
		if err != nil {
			return err
		}

		now := time.Now().UTC()
		monthAgo := now.AddDate(0, -1, 0)
		if err := u.listeningHistoryService.DeleteHistoriesByTimestamp(u.factory, user.ID, monthAgo); err != nil {
			return err
		}

		if err := client.ChangePlaylist(newReleases, user); err != nil {
			return err
		}
	}
	return nil
}
