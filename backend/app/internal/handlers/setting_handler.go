package handlers

import (
	"fmt"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/cookie"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/usecase"
	"github.com/labstack/echo"
	"log"
	"net/http"
)

func NewSettinHandler(
	logger *log.Logger,
	getSettingUsecase *usecase.GetSettingUsecase,
	editSettingUsecase *usecase.EditSettingUsecase) *SettingHandler {
	return &SettingHandler{
		logger:             logger,
		getSettingUsecase:  getSettingUsecase,
		editSettingUsecase: editSettingUsecase,
	}
}

type SettingHandler struct {
	logger             *log.Logger
	getSettingUsecase  *usecase.GetSettingUsecase
	editSettingUsecase *usecase.EditSettingUsecase
}

type UserPreference struct {
	IfRemixAdd    bool `json:"ifRemixAdd"`
	IfAcousticAdd bool `json:"ifAcousticAdd"`
}

func (h *SettingHandler) GetSettings(c echo.Context) error {
	fmt.Println("get settings")
	h.logger.Print("Get Settings")

	token, err := cookie.ReadCookie(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid cookie: %w", err)
	}

	ifRemixAdd, ifAcousticAdd, err := h.getSettingUsecase.GetSetting(token)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, UserPreference{IfRemixAdd: ifRemixAdd, IfAcousticAdd: ifAcousticAdd})
}

func (h *SettingHandler) EditSettings(c echo.Context) error {
	h.logger.Print("Edit Settings")

	token, err := cookie.ReadCookie(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid cookie: %w", err)
	}

	userPreference := new(UserPreference)
	if err := c.Bind(userPreference); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid parameters: %w", err)
	}

	if err := h.editSettingUsecase.EditSetting(token, userPreference.IfRemixAdd, userPreference.IfAcousticAdd); err != nil {
		return err
	}
	return nil
}
