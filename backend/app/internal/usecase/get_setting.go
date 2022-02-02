package usecase

import (
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/service"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/spotify_service"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/dao"
	"golang.org/x/oauth2"
)

func NewGetSettingUsecase(factory dao.Factory, config spotify_service.Config) *GetSettingUsecase {
	return &GetSettingUsecase{
		factory:                 factory,
		spotifyConfig:           config,
		userService:             service.NewUserService(),
		listeningHistoryService: service.NewListeningHistoryService(),
	}
}

type GetSettingUsecase struct {
	factory                 dao.Factory
	spotifyConfig           spotify_service.Config
	userService             *service.UserService
	listeningHistoryService *service.ListeningHistoryService
}

func (u *GetSettingUsecase) GetSetting(token *oauth2.Token) (bool, bool, error) {
	client := spotify_service.CreateNewClientByToken(u.spotifyConfig, token)

	userID, err := client.GetCurrentUserId()
	if err != nil {
		return false, false, err
	}

	user, err := u.userService.GetUser(u.factory, userID)
	if err != nil {
		return false, false, err
	}

	return user.IfRemixAdd, user.IfAcousticAdd, nil
}
