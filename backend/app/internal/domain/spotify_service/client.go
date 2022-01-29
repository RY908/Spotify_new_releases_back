package spotify_service

import (
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/entity"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

type Client struct {
	client *spotify.Client
}

func CreateSpotifyAccountServiceUrl(config Config) string {
	auth := spotify.NewAuthenticator(
		config.RedirectURI,
		spotify.ScopeUserReadRecentlyPlayed,
		spotify.ScopeUserReadPrivate,
		spotify.ScopePlaylistModifyPublic,
		spotify.ScopeUserFollowRead,
		spotify.ScopeImageUpload)
	auth.SetAuthInfo(config.ClientID, config.SecretKey)
	url := auth.AuthURL(config.State)

	return url
}

func CreateNewClientByUser(config Config, user entity.User) *Client {
	auth := spotify.NewAuthenticator(
		config.RedirectURI,
		spotify.ScopeUserReadRecentlyPlayed,
		spotify.ScopeUserReadPrivate,
		spotify.ScopePlaylistModifyPublic,
		spotify.ScopeUserFollowRead,
		spotify.ScopeImageUpload)
	auth.SetAuthInfo(config.ClientID, config.SecretKey)

	token := oauth2.Token{
		AccessToken:  user.AccessToken,
		TokenType:    user.TokenType,
		RefreshToken: user.RefreshToken,
		Expiry:       user.Expiry,
	}

	client := auth.NewClient(&token)

	return &Client{client: &client}
}

func CreateNewClientByToken(config Config, token oauth2.Token) *Client {
	auth := spotify.NewAuthenticator(
		config.RedirectURI,
		spotify.ScopeUserReadRecentlyPlayed,
		spotify.ScopeUserReadPrivate,
		spotify.ScopePlaylistModifyPublic,
		spotify.ScopeUserFollowRead,
		spotify.ScopeImageUpload)
	auth.SetAuthInfo(config.ClientID, config.SecretKey)
	client := auth.NewClient(&token)

	return &Client{client: &client}
}
