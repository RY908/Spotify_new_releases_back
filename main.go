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
	c.AddFunc("@every 30m", func() {
		UpdateRelation(mydbmap)
	})
	c.Start()
	defer c.Stop()

	r.HandleFunc("/", loginHandler)
	r.HandleFunc("/callback", redirectHandler)
	r.HandleFunc("/home", homeHandler).Methods("GET")
	//r.HandleFunc("/result", resultHander).Methods("POST")
	//.HandleFunc("/logout", logoutHandler)
	// rを割当
	//http.Handle("/", r)
	http.ListenAndServe(":8080", r)
}

