package handler

import (
	"fmt"
	"net/http"
	. "Spotify_new_releases/spotify"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("login")
	url := GetURL()
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers","Content-Type")
 	w.Header().Set("Access-Control-Allow-Credentials", "true")
	fmt.Println(url)
	http.Redirect(w, r, url, 302)
}