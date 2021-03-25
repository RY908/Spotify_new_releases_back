package main

import (
	"fmt"
	"net/http"
	"encoding/gob"
	"github.com/gorilla/mux"
	_ "github.com/go-sql-driver/mysql"
	. "Spotify_new_releases/database"
	// . "Spotify_new_releases/spotify"
	. "Spotify_new_releases/session"
	. "Spotify_new_releases/event"
	. "Spotify_new_releases/handler"
	"github.com/robfig/cron/v3"
)

// var (
// 	mydbmap = DatabaseInit()
// )

func main() {
	dbmap, err := DatabaseInit()
	if err != nil {
		fmt.Println(err)
	}
	// セッション初期処理
	gob.Register(UserSession{})
	SessionInit()

	r := mux.NewRouter()
	// url := GetURL()
	// fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	// cron
	c := cron.New()
	c.AddFunc("@every 20m", func() {
		if err := UpdateRelation(dbmap); err != nil {
			fmt.Println(err)
		}
	})
	c.AddFunc("10 00 * * 5", func() {
		if err := UpdatePlaylist(dbmap); err != nil {
			fmt.Println(err)
		}
	})
	c.AddFunc("@monthly", func() {
		if err := UpdateFollowingArtists(dbmap); err != nil {
			fmt.Println(err)
		}
	})
	c.Start()
	//defer c.Stop()

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		LoginHandler(w, r)
	}).Methods("GET")
	r.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		RedirectHandler(w, r, dbmap)
	})
	r.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		UserHandler(w, r, dbmap)
	}).Methods("GET")
	r.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		DeleteHandler(w, r, dbmap)
	}).Methods("POST")
	r.HandleFunc("/setting", func(w http.ResponseWriter, r *http.Request) {
		SettingHandler(w, r, dbmap)
	}).Methods("GET")
	r.HandleFunc("/setting/save", func(w http.ResponseWriter, r *http.Request) {
		SettingEditHandler(w, r, dbmap)
	}).Methods("POST")
	// r.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
	// 	UpdatePlaylist(mydbmap)
	// }).Methods("GET")

	http.ListenAndServe(":9990", r)
}

