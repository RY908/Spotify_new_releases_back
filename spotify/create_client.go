package spotify

import (
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"net/http"
	"os"
	"fmt"
	"Spotify_new_releases/database"
)

const redirectURI = "http://localhost:8080/callback"

var (
	clientID = os.Getenv("SPOTIFY_ID_3")
	secretKey = os.Getenv("SPOTIFY_SECRET_3")
	state = "abc123"
)

type Client struct {
	Client *spotify.Client
}

func GetURL() string {
	auth  := spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadRecentlyPlayed, spotify.ScopeUserReadPrivate, spotify.ScopePlaylistModifyPublic)
	auth.SetAuthInfo(clientID, secretKey)
	url := auth.AuthURL(state)

	return url
}

func CreateMyClient(r *http.Request) (*Client, *oauth2.Token, *http.Request) {
	auth  := spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadRecentlyPlayed, spotify.ScopeUserReadPrivate, spotify.ScopePlaylistModifyPublic)
	auth.SetAuthInfo(clientID, secretKey)
	token, err := auth.Token(state, r)
	if err != nil {
		fmt.Println(err)
	}
	client := auth.NewClient(token)

	return &Client{Client: &client}, token, r
}

func CreateMyClientFromUserInfo(user database.UserInfo) *Client {
	auth  := spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadRecentlyPlayed, spotify.ScopeUserReadPrivate, spotify.ScopePlaylistModifyPublic)
	auth.SetAuthInfo(clientID, secretKey)

	accessToken := user.AccessToken
	tokenType := user.TokenType
	refreshToken := user.RefreshToken
	expiry := user.Expiry

	// make token
	token := oauth2.Token{AccessToken:accessToken, TokenType:tokenType, RefreshToken:refreshToken, Expiry:expiry}

	// create new client
	client := auth.NewClient(&token)

	return &Client{Client: &client}
}

func CreateMyClientFromToken(token oauth2.Token) *Client {
	auth  := spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadRecentlyPlayed, spotify.ScopeUserReadPrivate, spotify.ScopePlaylistModifyPublic)
	auth.SetAuthInfo(clientID, secretKey)
	client := auth.NewClient(&token)

	return &Client{Client: &client}
}