package cookie

import (
	"fmt"
	"net/http"
	"os"
	"time"

	. "github.com/RY908/Spotify_new_releases_back/backend/database"
	. "github.com/RY908/Spotify_new_releases_back/backend/spotify"

	"golang.org/x/oauth2"
)

var (
	layout = "2006-01-02 15:04:05"
	DOMAIN = os.Getenv("DOMAIN")
)

// set oauth token in cookie
func SetCookie(w http.ResponseWriter, token *oauth2.Token) (http.ResponseWriter, error) {
	http.SetCookie(w, &http.Cookie{
		Name:   "token",
		Value:  token.AccessToken,
		Path:   "/",
		Domain: DOMAIN,
	})

	http.SetCookie(w, &http.Cookie{
		Name:   "tokenType",
		Value:  token.TokenType,
		Path:   "/",
		Domain: DOMAIN,
	})

	http.SetCookie(w, &http.Cookie{
		Name:   "refreshToken",
		Value:  token.RefreshToken,
		Path:   "/",
		Domain: DOMAIN,
	})

	http.SetCookie(w, &http.Cookie{
		Name:   "expiry",
		Value:  (token.Expiry).Format(layout),
		Path:   "/",
		Domain: DOMAIN,
	})
	return w, nil
}

// get access token from cookie
func GetToken(r *http.Request) (oauth2.Token, error) {
	tokenCookie, err := r.Cookie("token")
	if err != nil {
		err = fmt.Errorf("unable to get token from cookie: %w", err)
		return oauth2.Token{}, err
	}
	refreshTokenCookie, err := r.Cookie("refreshToken")
	if err != nil {
		err = fmt.Errorf("unable to get refresh token from cookie: %w", err)
		return oauth2.Token{}, err
	}
	tokenTypeCookie, err := r.Cookie("tokenType")
	if err != nil {
		err = fmt.Errorf("unable to get token type from cookie: %w", err)
		return oauth2.Token{}, err
	}
	expiryCookie, err := r.Cookie("expiry")
	if err != nil {
		err = fmt.Errorf("unable to get expiry from cookie: %w", err)
		return oauth2.Token{}, err
	}
	accessToken := tokenCookie.Value
	refreshToken := refreshTokenCookie.Value
	tokenType := tokenTypeCookie.Value
	expiryString := expiryCookie.Value
	expiry, _ := time.Parse(layout, expiryString)

	// get token, client and user id
	token := oauth2.Token{AccessToken: accessToken, TokenType: tokenType, RefreshToken: refreshToken, Expiry: expiry}
	return token, nil
}

// get access token and create client then check if the client is already in database
func GetUser(r *http.Request, mydbmap *MyDbMap, token oauth2.Token) (bool, UserInfo, error) {
	// token, err := GetToken(r)
	// if err != nil {
	// 	return false, UserInfo{}, err
	// }
	client := CreateMyClientFromToken(token)
	userId, err := client.GetCurrentUserId()
	if err != nil {
		err = fmt.Errorf("unable to get userId: %w", err)
		return false, UserInfo{}, err
	}
	exists, user, err := mydbmap.ExistUser(userId)
	if err != nil {
		err = fmt.Errorf("unable to check if the user exists in database: %w", err)
		return false, UserInfo{}, err
	}
	return exists, user, nil
}
