package handlers

import (
	"fmt"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/cookie"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/spotify_service"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/usecase"
	"github.com/labstack/echo"
	"net/http"
	"os"
)

var (
	errURI = os.Getenv("ERR_URI")
	sucURI = os.Getenv("SUC_URI")
)

func NewLoginHandler(createPlaylistUsecase *usecase.CreatePlaylistUsecase) *LoginHandler {
	return &LoginHandler{
		createPlaylistUsecase: createPlaylistUsecase,
	}
}

type LoginHandler struct {
	createPlaylistUsecase *usecase.CreatePlaylistUsecase
}

func (h *LoginHandler) Login(c echo.Context) error {
	url := spotify_service.GetURL()
	fmt.Println(url)
	c.Redirect(http.StatusFound, url)
	return nil
}

func (h *LoginHandler) Callback(c echo.Context) error {
	token, err := spotify_service.GetToken(c.Request())
	fmt.Println(token, err)
	if err != nil {
		return err
	}

	cookie.WriteCookie(c, token)

	if err := h.createPlaylistUsecase.CreatePlaylist(token); err != nil {
		fmt.Println(err)
		return err
	}

	c.Redirect(http.StatusFound, sucURI+"/"+token.AccessToken)
	return nil
}
