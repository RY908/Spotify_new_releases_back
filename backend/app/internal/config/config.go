package config

import (
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/spotify_service"
	"os"
)

type Config struct {
	DBConfig       *DBConfig
	SpotifyConfig  *spotify_service.Config
	CallbackConfig *CallbackConfig
	APIConfig      *APIConfig
}

type DBConfig struct {
	DBPath     string
	TestDBPath string
}

type CallbackConfig struct {
	SuccessURI string
	ErrorURI   string
}

type APIConfig struct {
	AccessControlAllowOrigin string
}

func LoadConfig() *Config {
	return &Config{
		DBConfig: &DBConfig{
			DBPath:     os.Getenv("SQL_PATH"),
			TestDBPath: os.Getenv("SQL_PATH_TEST"),
		},
		SpotifyConfig: &spotify_service.Config{
			RedirectURI: os.Getenv("REDIRECT_URI"),
			ClientID:    os.Getenv("SPOTIFY_ID_3"),
			SecretKey:   os.Getenv("SPOTIFY_SECRET_3"),
			State:       "abc123",
		},
		CallbackConfig: &CallbackConfig{
			SuccessURI: os.Getenv("SUC_URI"),
			ErrorURI:   os.Getenv("ERR_URI"),
		},
		APIConfig: &APIConfig{
			AccessControlAllowOrigin: os.Getenv("ACCESS_CONTROL_ALLOW_ORIGIN"),
		},
	}
}
