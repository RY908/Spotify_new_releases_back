package cookie 

import (
	"fmt"
	"time"
	"net/http"
	"golang.org/x/oauth2"
	. "Spotify_new_releases/spotify"
	. "Spotify_new_releases/database"
)

var layout = "2006-01-02 15:04:05"

func SetCookie(w http.ResponseWriter, token *oauth2.Token) (http.ResponseWriter, error) {
	http.SetCookie(w, &http.Cookie{
        Name: "token",
        Value: token.AccessToken,
     	Path: "/",
	})
	http.SetCookie(w, &http.Cookie{
		Name: "tokenType",
		Value: token.TokenType,
		Path: "/",
	})
	http.SetCookie(w, &http.Cookie{
		Name: "refreshToken",
		Value: token.RefreshToken,
		Path: "/",
	})
	http.SetCookie(w, &http.Cookie{
		Name: "expiry",
		Value: (token.Expiry).Format(layout),
		Path: "/",
	})
	return w, nil
}

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
	token := oauth2.Token{AccessToken:accessToken, TokenType:tokenType, RefreshToken:refreshToken, Expiry:expiry}
	return token, nil
}

func GetUser(r *http.Request, mydbmap *MyDbMap) (bool, UserInfo, error) {
	token, err := GetToken(r)
	if err != nil {
		fmt.Println(err)
	}
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