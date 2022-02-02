package handlers

import (
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/entity"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/spotify_service"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/usecase"
	"github.com/labstack/echo"
	"net/http"
)

type userHandler struct {
	deleteListeningHistoryUsecase *usecase.DeleteListeningHistoryUsecase
	getArtistsByUserIDUsecase     *usecase.GetArtistsByUserIDUsecase
}

type Artists struct {
	IDs []string `json:"artistsId"`
}

type UserArtists struct {
	Artists []*entity.UserArtist `json:"artists"`
}

func (h *userHandler) DeleteArtists(c echo.Context) error {
	token, err := spotify_service.GetToken(c.Request())
	if err != nil {
		return err
	}

	// TODO: cookie

	artistsReq := new(Artists)
	if err := c.Bind(artistsReq); err != nil {
		return err
	}

	if err := h.deleteListeningHistoryUsecase.DeleteListeningHistory(token, artistsReq.IDs); err != nil {
		return err
	}

	artists, err := h.getArtistsByUserIDUsecase.GetArtistsByUserIDUsecase(token)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, UserArtists{Artists: artists})
}

func (h *userHandler) GetArtists(c echo.Context) error {
	token, err := spotify_service.GetToken(c.Request())
	if err != nil {
		return err
	}

	// TODO: cookie

	artists, err := h.getArtistsByUserIDUsecase.GetArtistsByUserIDUsecase(token)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, UserArtists{Artists: artists})
}
