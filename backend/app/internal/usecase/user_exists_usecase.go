package usecase

import (
	"database/sql"
	"errors"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/service"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/spotify_service"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/dao"
	"golang.org/x/oauth2"
	"log"
)

func NewUserExistsUsecase(factory dao.Factory, logger *log.Logger, config *spotify_service.Config) *UserExistsUsecase {
	return &UserExistsUsecase{
		factory:       factory,
		logger:        logger,
		spotifyConfig: config,
		userService:   service.NewUserService(),
	}
}

type UserExistsUsecase struct {
	factory       dao.Factory
	logger        *log.Logger
	spotifyConfig *spotify_service.Config
	userService   *service.UserService
}

func (u *UserExistsUsecase) IsUserExists(token *oauth2.Token) (bool, error) {
	client := spotify_service.CreateNewClientByToken(u.spotifyConfig, token)

	userID, err := client.GetCurrentUserId()
	if err != nil {
		u.logger.Print(err)
		return false, err
	}

	_, err = u.userService.GetUser(u.factory, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		} else {
			return false, nil
		}
	}
	return true, nil
}
