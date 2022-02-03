package config

import (
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/spotify_service"
	"os"
)

type Config struct {
	DBConfig      DBConfig
	SpotifyConfig spotify_service.Config
}

type DBConfig struct {
	DBPath     string
	TestDBPath string
}

func LoadConfig() Config {
	return Config{
		DBConfig: DBConfig{
			DBPath:     os.Getenv("SQL_PATH"),
			TestDBPath: os.Getenv("SQL_PATH_TEST"),
		},
		SpotifyConfig: spotify_service.Config{
			RedirectURI: os.Getenv("REDIRECT_URI"),
			ClientID:    os.Getenv("SPOTIFY_ID_3"),
			SecretKey:   os.Getenv("SPOTIFY_SECRET_3"),
			State:       "abc123",
		},
	}
}
