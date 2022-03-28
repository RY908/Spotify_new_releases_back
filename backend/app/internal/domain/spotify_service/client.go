package spotify_service

import (
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/entity"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"net/http"
)

type Client struct {
	client *spotify.Client
}

func GetURL(config *Config) string {
	auth := newAuthenticator(config.RedirectURI, config.ClientID, config.SecretKey)
	url := auth.AuthURL(config.State)

	return url
}

func GetToken(config *Config, r *http.Request) (*oauth2.Token, error) {
	auth := newAuthenticator(config.RedirectURI, config.ClientID, config.SecretKey)
	token, err := auth.Token(config.State, r)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func CreateNewClientByUser(config *Config, user *entity.User) *Client {
	auth := newAuthenticator(config.RedirectURI, config.ClientID, config.SecretKey)

	token := oauth2.Token{
		AccessToken:  user.AccessToken,
		TokenType:    user.TokenType,
		RefreshToken: user.RefreshToken,
		Expiry:       user.Expiry,
	}

	client := auth.NewClient(&token)

	return &Client{client: &client}
}

func CreateNewClientByToken(config *Config, token *oauth2.Token) *Client {
	auth := newAuthenticator(config.RedirectURI, config.ClientID, config.SecretKey)
	client := auth.NewClient(token)

	return &Client{client: &client}
}

func newAuthenticator(redirectURI, clientID, secretKey string) spotify.Authenticator {
	auth := spotify.NewAuthenticator(
		redirectURI,
		spotify.ScopeUserReadRecentlyPlayed,
		spotify.ScopeUserReadPrivate,
		spotify.ScopePlaylistModifyPublic,
		spotify.ScopeUserFollowRead,
		spotify.ScopeImageUpload)
	auth.SetAuthInfo(clientID, secretKey)
	return auth
}
