package main

import (
	"fmt"
	"log"
	"os"
	"net/http"
	"encoding/gob"
	"golang.org/x/oauth2"
	"github.com/zmb3/spotify"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-gorp/gorp"
)

const redirectURI = "http://localhost:8080/callback"

var (
	clientID = os.Getenv("SPOTIFY_ID_3")
	secretKey = os.Getenv("SPOTIFY_SECRET_3")
	auth  = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadRecentlyPlayed, spotify.ScopeUserReadPrivate)
	state = "abc123"
	session_name = "spotify_access_token"
	store *sessions.CookieStore
	session *sessions.Session
	sqlPath = os.Getenv("SQL_PATH")
	db, _ = sql.Open("mysql", sqlPath)
	dbmap = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{}}
)

type UserSession struct {
	ID 		string
	Token 	oauth2.Token
}

func main() {
	dbmap.AddTableWithName(ArtistInfo{}, "Artist").SetKeys(false, "ArtistId")
	dbmap.AddTableWithName(ListenTo{}, "ListenTo").SetKeys(true, "ListenId")
	dbmap.AddTableWithName(UserInfo{}, "User").SetKeys(false, "UserId")
	//dbmap.CreateTablesIfNotExists()
	fmt.Printf("dbmap: %T", dbmap)
	defer db.Close()
	defer dbmap.Db.Close()

	auth.SetAuthInfo(clientID, secretKey)

	gob.Register(UserSession{})

	// セッション初期処理
	sessionInit()

	r := mux.NewRouter()
	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	r.HandleFunc("/callback", redirectHandler)
	r.HandleFunc("/home", homeHandler).Methods("GET")
	//r.HandleFunc("/result", resultHander).Methods("POST")
	//.HandleFunc("/logout", logoutHandler)
	// rを割当
	//http.Handle("/", r)
	http.ListenAndServe(":8080", r)
}

