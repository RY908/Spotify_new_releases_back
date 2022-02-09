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

func NewLoginHandler(logger *log.Logger, config *config.CallbackConfig, createPlaylistUsecase *usecase.CreatePlaylistUsecase) *LoginHandler {
	return &LoginHandler{
		logger:                logger,
		config:                config,
		createPlaylistUsecase: createPlaylistUsecase,
	}
}

type LoginHandler struct {
	logger                *log.Logger
	config                *config.CallbackConfig
	createPlaylistUsecase *usecase.CreatePlaylistUsecase
}

func (h *LoginHandler) Login(c echo.Context) error {
	h.logger.Print("Login")
	url := spotify_service.GetURL()
	c.Redirect(http.StatusFound, url)
	return nil
}

func (h *LoginHandler) Callback(c echo.Context) error {
	h.logger.Print("Callback")
	token, err := spotify_service.GetToken(c.Request())
	if err != nil {
		return err
	}

	cookie.WriteCookie(c, token)

	if err := h.createPlaylistUsecase.CreatePlaylist(token); err != nil {
		return err
	}

	c.Redirect(http.StatusFound, h.config.SuccessURI+"/"+token.AccessToken)
	return nil
}
