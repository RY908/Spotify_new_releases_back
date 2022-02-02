package usecase

import (
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/service"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/spotify_service"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/dao"
	"time"
)

func NewUpdateFollowingArtistsUsecase(factory dao.Factory, config spotify_service.Config) *UpdateFollowingArtistsUsecase {
	return &UpdateFollowingArtistsUsecase{
		factory:                 factory,
		spotifyConfig:           config,
		artistService:           service.NewArtistService(),
		userService:             service.NewUserService(),
		listeningHistoryService: service.NewListeningHistoryService(),
	}
}

type UpdateFollowingArtistsUsecase struct {
	factory                 dao.Factory
	spotifyConfig           spotify_service.Config
	artistService           *service.ArtistService
	userService             *service.UserService
	listeningHistoryService *service.ListeningHistoryService
}

func (u *UpdateFollowingArtistsUsecase) UpdateFollowingArtists() error {
	users, err := u.userService.GetAllUsers(u.factory)
	if err != nil {
		return err
	}
	for _, user := range users {
		client := spotify_service.CreateNewClientByUser(u.spotifyConfig, *user)

		artists, err := client.GetFollowingArtists(user.ID)
		if err != nil {
			return err
		}

		if err := u.artistService.InsertArtists(u.factory, artists); err != nil {
			return err
		}

		timestamp := time.Now()
		if err := u.listeningHistoryService.UpdateIsFollowings(u.factory, artists, user.ID, timestamp); err != nil {
			return err
		}

		if err := u.listeningHistoryService.DeleteHistoriesByTimestamp(u.factory, user.ID, timestamp); err != nil {
			return err
		}
	}
	return nil
}
