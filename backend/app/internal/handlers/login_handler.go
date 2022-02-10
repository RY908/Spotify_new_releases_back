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
	createPlaylistUsecase *usecase.CreatePlaylistUsecase) *LoginHandler {
	return &LoginHandler{
		logger:                logger,
		callbackConfig:        callbackConfig,
		spotifyConfig:         spotifyConfig,
		createPlaylistUsecase: createPlaylistUsecase,
	}
}

type LoginHandler struct {
	logger                *log.Logger
	callbackConfig        *config.CallbackConfig
	spotifyConfig         *spotify_service.Config
	createPlaylistUsecase *usecase.CreatePlaylistUsecase
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
		return err
	}

	// TODO: redirect to login page
	cookie.WriteCookie(c, token)

	if err := h.createPlaylistUsecase.CreatePlaylist(token); err != nil {
		return err
	}

	c.Redirect(http.StatusFound, h.callbackConfig.SuccessURI+"/"+token.AccessToken)
	return nil
}
