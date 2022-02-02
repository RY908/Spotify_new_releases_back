package handlers

import (
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/spotify_service"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/usecase"
	"github.com/labstack/echo"
	"net/http"
)

type SettingHandler struct {
	getSettingUsecase  *usecase.GetSettingUsecase
	editSettingUsecase *usecase.EditSettingUsecase
}

type UserPreference struct {
	IfRemixAdd    bool `json:"ifRemixAdd"`
	IfAcousticAdd bool `json:"ifAcousticAdd"`
}

func (h *SettingHandler) GetSettings(c echo.Context) error {
	token, err := spotify_service.GetToken(c.Request())
	if err != nil {
		return err
	}

	// TODO: cookie

	ifRemixAdd, ifAcousticAdd, err := h.getSettingUsecase.GetSetting(token)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, UserPreference{IfRemixAdd: ifRemixAdd, IfAcousticAdd: ifAcousticAdd})
}

func (h *SettingHandler) EditSettings(c echo.Context) error {
	token, err := spotify_service.GetToken(c.Request())
	if err != nil {
		return err
	}

	// TODO: cookie

	userPreference := new(UserPreference)
	if err := c.Bind(userPreference); err != nil {
		return err
	}

	if err := h.editSettingUsecase.EditSetting(token, userPreference.IfRemixAdd, userPreference.IfAcousticAdd); err != nil {
		return err
	}
	return nil
}
