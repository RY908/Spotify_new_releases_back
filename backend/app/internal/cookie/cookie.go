package cookie

import (
	"github.com/labstack/echo"
	"golang.org/x/oauth2"
	"net/http"
	"time"
)

var (
	layout = "2006-01-02 15:04:05"
	//DOMAIN = os.Getenv("DOMAIN")
)

func WriteCookie(c echo.Context, token *oauth2.Token) {
	c.SetCookie(&http.Cookie{
		Name:  "token",
		Value: token.AccessToken,
		Path:  "/",
		//Domain: DOMAIN,
	})
	c.SetCookie(&http.Cookie{
		Name:  "tokenType",
		Value: token.TokenType,
		Path:  "/",
		//Domain: DOMAIN,
	})
	c.SetCookie(&http.Cookie{
		Name:  "refreshToken",
		Value: token.RefreshToken,
		Path:  "/",
		//Domain: DOMAIN,
	})
	c.SetCookie(&http.Cookie{
		Name:  "expiry",
		Value: (token.Expiry).Format(layout),
		Path:  "/",
		//Domain: DOMAIN,
	})
}

func ReadCookie(c echo.Context) (*oauth2.Token, error) {
	tokenCookie, err := c.Cookie("token")
	if err != nil {
		return nil, err
	}
	tokenTypeCookie, err := c.Cookie("tokenType")
	if err != nil {
		return nil, err
	}
	refreshTokenCookie, err := c.Cookie("refreshToken")
	if err != nil {
		return nil, err
	}
	expiryCookie, err := c.Cookie("expiry")
	if err != nil {
		return nil, err
	}
	expiryString := expiryCookie.Value
	expiry, _ := time.Parse(layout, expiryString)

	return &oauth2.Token{
		AccessToken:  tokenCookie.Value,
		TokenType:    tokenTypeCookie.Value,
		RefreshToken: refreshTokenCookie.Value,
		Expiry:       expiry,
	}, nil
}
