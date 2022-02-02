package spotify_service

import (
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/entity"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"net/http"
	"os"
)

var (
	redirectURI = os.Getenv("REDIRECT_URI")
	clientID    = os.Getenv("SPOTIFY_ID_3")
	secretKey   = os.Getenv("SPOTIFY_SECRET_3")
	state       = "abc123"
)

type Client struct {
	client *spotify.Client
}

func GetURL() string {
	auth := spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadRecentlyPlayed, spotify.ScopeUserReadPrivate,
		spotify.ScopePlaylistModifyPublic, spotify.ScopeUserFollowRead, spotify.ScopeImageUpload)
	auth.SetAuthInfo(clientID, secretKey)
	url := auth.AuthURL(state)

	return url
}

func GetToken(r *http.Request) (*oauth2.Token, error) {
	auth := spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadRecentlyPlayed, spotify.ScopeUserReadPrivate,
		spotify.ScopePlaylistModifyPublic, spotify.ScopeUserFollowRead, spotify.ScopeImageUpload)
	auth.SetAuthInfo(clientID, secretKey)
	token, err := auth.Token(state, r)
	if err != nil {
		return nil, err
	}
	return token, nil
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

func CreateNewClientByToken(config Config, token *oauth2.Token) *Client {
	auth := spotify.NewAuthenticator(
		config.RedirectURI,
		spotify.ScopeUserReadRecentlyPlayed,
		spotify.ScopeUserReadPrivate,
		spotify.ScopePlaylistModifyPublic,
		spotify.ScopeUserFollowRead,
		spotify.ScopeImageUpload)
	auth.SetAuthInfo(config.ClientID, config.SecretKey)
	client := auth.NewClient(token)

	return &Client{client: &client}
}
