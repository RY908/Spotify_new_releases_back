package spotify

import (
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"net/http"
	"os"
	"fmt"
	"Spotify_new_releases/database"
)

var (
	redirectURI = "https://newreleases.tk/api/callback" //os.Getenv("REDIRECT_URI")
	clientID = os.Getenv("SPOTIFY_ID_3")
	secretKey = os.Getenv("SPOTIFY_SECRET_3")
	state = "abc123"
)

type Client struct {
	Client *spotify.Client
}

func GetURL() string {
	fmt.Println(redirectURI)
	auth  := spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadRecentlyPlayed, spotify.ScopeUserReadPrivate, spotify.ScopePlaylistModifyPublic, spotify.ScopeUserFollowRead, spotify.ScopeImageUpload)
	auth.SetAuthInfo(clientID, secretKey)
	url := auth.AuthURL(state)

	return url
}

// CreateMyClient creates a new client.
// This is called when the user first logs in.
func CreateMyClient(r *http.Request) (*Client, *oauth2.Token, *http.Request, error) {
	auth  := spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadRecentlyPlayed, spotify.ScopeUserReadPrivate, spotify.ScopePlaylistModifyPublic, spotify.ScopeUserFollowRead, spotify.ScopeImageUpload)
	auth.SetAuthInfo(clientID, secretKey)
	token, err := auth.Token(state, r)
	if err != nil {
		err = fmt.Errorf("unable to get token: %w", err)
		return nil, nil, r, err
	}
	client := auth.NewClient(token)

	return &Client{Client: &client}, token, r, nil
}

// CreateMyClientFromUserInfo creates a new client from data in the database.
func CreateMyClientFromUserInfo(user database.UserInfo) *Client {
	auth  := spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadRecentlyPlayed, spotify.ScopeUserReadPrivate, spotify.ScopePlaylistModifyPublic, spotify.ScopeUserFollowRead, spotify.ScopeImageUpload)
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

// CreateMyClientFromToken creates a new client from oauth2 token.
func CreateMyClientFromToken(token oauth2.Token) *Client {
	auth  := spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadRecentlyPlayed, spotify.ScopeUserReadPrivate, spotify.ScopePlaylistModifyPublic, spotify.ScopeUserFollowRead, spotify.ScopeImageUpload)
	auth.SetAuthInfo(clientID, secretKey)
	client := auth.NewClient(&token)

	return &Client{Client: &client}
}

func CreateMyClientFromCode(r *http.Request) (*Client, *oauth2.Token, *http.Request, error) {
	auth  := spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadRecentlyPlayed, spotify.ScopeUserReadPrivate, spotify.ScopePlaylistModifyPublic, spotify.ScopeUserFollowRead, spotify.ScopeImageUpload)
	code := r.FormValue("code")
	token, err := auth.Exchange(code)
	if err != nil {
		err = fmt.Errorf("unable to get token from code: %w", err)
		return nil, nil, r, err
	}
	client := auth.NewClient(token)

	return &Client{Client: &client}, token, r, nil
}
