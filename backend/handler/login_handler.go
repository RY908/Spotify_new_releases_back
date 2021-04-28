package handler

import (
	. "Spotify_new_releases/backend/spotify"
	"fmt"
	"net/http"
	"os"
)

var (
	accessControlAllowOrigin = os.Getenv("ACCESS_CONTROL_ALLOW_ORIGIN")
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("login")
	url := GetURL()
	w.Header().Set("Access-Control-Allow-Origin", accessControlAllowOrigin)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	fmt.Println(url)
	http.Redirect(w, r, url, 302)
}
