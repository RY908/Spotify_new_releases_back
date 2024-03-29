package handlers

import (
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/cookie"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/entity"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/usecase"
	"github.com/labstack/echo"
	"log"
	"net/http"
)

func NewUserHandler(
	logger *log.Logger,
	deleteListeningHistoryUsecase *usecase.DeleteListeningHistoryUsecase,
	getArtistsByUserIDUsecase *usecase.GetArtistsByUserIDUsecase) *UserHandler {
	return &UserHandler{
		logger:                        logger,
		deleteListeningHistoryUsecase: deleteListeningHistoryUsecase,
		getArtistsByUserIDUsecase:     getArtistsByUserIDUsecase,
	}
}

type UserHandler struct {
	logger                        *log.Logger
	deleteListeningHistoryUsecase *usecase.DeleteListeningHistoryUsecase
	getArtistsByUserIDUsecase     *usecase.GetArtistsByUserIDUsecase
}

type Artists struct {
	IDs []string `json:"artistsId"`
}

type UserArtists struct {
	Artists []*entity.UserArtist `json:"artists"`
}

func (h *UserHandler) DeleteArtists(c echo.Context) error {
	h.logger.Print("Delete Artists")

	token, err := cookie.ReadCookie(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid cookie: %w", err)
	}

	artistsReq := new(Artists)
	if err := c.Bind(artistsReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid parameters: %w", err)
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

func (h *UserHandler) GetArtists(c echo.Context) error {
	h.logger.Print("Get Artists")

	token, err := cookie.ReadCookie(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid cookie: %w", err)
	}

	artists, err := h.getArtistsByUserIDUsecase.GetArtistsByUserIDUsecase(token)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, UserArtists{Artists: artists})
}
