package handlers

import (
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/spotify_service"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/usecase"
	"github.com/labstack/echo"
	"net/http"
)

type LoginHandler struct {
	createPlaylistUsecase *usecase.CreatePlaylistUsecase
}

func (h *LoginHandler) Login(c echo.Context) {
	url := spotify_service.GetURL()
	c.Redirect(http.StatusFound, url)
}

func (h *LoginHandler) Callback(c echo.Context) error {
	token, err := spotify_service.GetToken(c.Request())
	if err != nil {
		return err
	}

	// TODO: cookie

	if err := h.createPlaylistUsecase.CreatePlaylist(token); err != nil {
		return err
	}
	return nil
}
