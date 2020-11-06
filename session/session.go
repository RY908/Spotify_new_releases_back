package session

import (
	"net/http"
	"golang.org/x/oauth2"
	"github.com/gorilla/sessions"
	"crypto/rand"
	"encoding/base32"
	"io"
	"strings"
	"fmt"
)

var (
	Session_name = "spotify_access_token"
	Store *sessions.CookieStore
	Session *sessions.Session
)

type UserSession struct {
	ID 			string
	Token 		oauth2.Token
	PlaylistId	string
}

func GetTokenFromSession(r *http.Request) (string, oauth2.Token, string) {
	session, _ := Store.Get(r, Session_name)
	userId :=session.Values["user"].(UserSession).ID
	token := session.Values["user"].(UserSession).Token
	playlistId := session.Values["user"].(UserSession).PlaylistId
	return userId, token, playlistId
}

func SessionInit(){

	// 乱数生成
	b := make([]byte, 48)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		panic(err)
	}
	str := strings.TrimRight(base32.StdEncoding.EncodeToString(b), "=")

	// 新しいstoreとセッションを準備
	Store = sessions.NewCookieStore([]byte(str))
	Session = sessions.NewSession(Store, Session_name)

	// セッションの有効範囲を指定
	Store.Options = &sessions.Options{
		Domain:     "localhost",
		Path:       "/",
		MaxAge:     0, // the cookie will be deleted after the browser session ends.
		Secure:     false,
		HttpOnly:   true,
	}

	// log
	fmt.Println("key     data --")
	fmt.Println(str)
	fmt.Println("")
	fmt.Println("store   data --")
	fmt.Println(Store)
	fmt.Println("")
	fmt.Println("session data --")
	fmt.Println(Session)
	fmt.Println("")

}