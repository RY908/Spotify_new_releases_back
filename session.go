package main

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

func getTokenFromSession(r *http.Request) (string, oauth2.Token) {
	session, _ := store.Get(r, session_name)
	userId :=session.Values["user"].(UserSession).ID
	token := session.Values["user"].(UserSession).Token
	return userId, token
}

func sessionInit(){

	// 乱数生成
	b := make([]byte, 48)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		panic(err)
	}
	str := strings.TrimRight(base32.StdEncoding.EncodeToString(b), "=")

	// 新しいstoreとセッションを準備
	store = sessions.NewCookieStore([]byte(str))
	session = sessions.NewSession(store, session_name)

	// セッションの有効範囲を指定
	store.Options = &sessions.Options{
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
	fmt.Println(store)
	fmt.Println("")
	fmt.Println("session data --")
	fmt.Println(session)
	fmt.Println("")

}