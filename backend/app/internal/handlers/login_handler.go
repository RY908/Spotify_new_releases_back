package handlers

import (
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/config"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/cookie"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/spotify_service"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/usecase"
	"github.com/labstack/echo"
	"log"
	"net/http"
)

func NewLoginHandler(
	logger *log.Logger,
	callbackConfig *config.CallbackConfig,
	spotifyConfig *spotify_service.Config,
	initializePlaylistUsecase *usecase.InitializePlaylistUsecase,
	initializeFollowingUsecase *usecase.InitializeFollowingUsecase,
	userExistsUsecase *usecase.UserExistsUsecase) *LoginHandler {
	return &LoginHandler{
		logger:                     logger,
		callbackConfig:             callbackConfig,
		spotifyConfig:              spotifyConfig,
		initializePlaylistUsecase:  initializePlaylistUsecase,
		initializeFollowingUsecase: initializeFollowingUsecase,
		userExistsUsecase:          userExistsUsecase,
	}
}

type LoginHandler struct {
	logger                     *log.Logger
	callbackConfig             *config.CallbackConfig
	spotifyConfig              *spotify_service.Config
	initializePlaylistUsecase  *usecase.InitializePlaylistUsecase
	initializeFollowingUsecase *usecase.InitializeFollowingUsecase
	userExistsUsecase          *usecase.UserExistsUsecase
}

func (h *LoginHandler) Login(c echo.Context) error {
	h.logger.Print("Login")
	url := spotify_service.GetURL(h.spotifyConfig)
	c.Redirect(http.StatusFound, url)
	return nil
}

func (h *LoginHandler) Callback(c echo.Context) error {
	h.logger.Print("Callback")
	token, err := spotify_service.GetToken(h.spotifyConfig, c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid cookie: %w", err)
	}

	cookie.WriteCookie(c, token)

	var ok = false
	if ok, err = h.userExistsUsecase.IsUserExists(token); err != nil {
		return err
	}

	if !ok {
		if err := h.initializePlaylistUsecase.CreatePlaylist(token); err != nil {
			return c.Redirect(http.StatusFound, h.callbackConfig.ErrorURI)
		}
		if err := h.initializeFollowingUsecase.InsertFollowingUsecase(token); err != nil {
			return err
		}
	}

	return c.Redirect(http.StatusFound, h.callbackConfig.SuccessURI+"/"+token.AccessToken)
}
