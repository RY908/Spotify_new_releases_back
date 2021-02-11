package main

import (
	"fmt"
	"net/http"
	"encoding/gob"
	"github.com/gorilla/mux"
	_ "github.com/go-sql-driver/mysql"
	. "Spotify_new_releases/database"
	. "Spotify_new_releases/spotify"
	. "Spotify_new_releases/session"
	. "Spotify_new_releases/event"
	. "Spotify_new_releases/handler"
	"github.com/robfig/cron/v3"
)

var (
	mydbmap = DatabaseInit()
)

func main() {
	// セッション初期処理
	gob.Register(UserSession{})
	SessionInit()

	r := mux.NewRouter()
	url := GetURL()
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	// cron
	c := cron.New()
	c.AddFunc("@every 20m", func() {
		if err := UpdateRelation(mydbmap); err != nil {
			fmt.Println(err)
		}
	})
	c.AddFunc("10 00 * * 5", func() {
		if err := UpdatePlaylist(mydbmap); err != nil {
			fmt.Println(err)
		}
	})
	c.Start()
	//defer c.Stop()

	r.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		LoginHandler(w, r)
	}).Methods("GET")
	r.HandleFunc("/api/callback", func(w http.ResponseWriter, r *http.Request) {
		RedirectHandler(w, r, mydbmap)
	})
	r.HandleFunc("/api/user", func(w http.ResponseWriter, r *http.Request) {
		UserHandler(w, r, mydbmap)
	}).Methods("GET")
	r.HandleFunc("/api/delete", func(w http.ResponseWriter, r *http.Request) {
		DeleteHandler(w, r, mydbmap)
	}).Methods("POST")

	http.ListenAndServe(":9990", r)
}

