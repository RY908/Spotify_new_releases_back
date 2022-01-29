package config

import "github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/spotify_service"

type Config struct {
	SpotifyConfig spotify_service.Config
}
