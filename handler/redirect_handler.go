package handler

import (
	"fmt"
	//"os"
	"net/http"
	. "Spotify_new_releases/spotify"
	//. "Spotify_new_releases/session"
	. "Spotify_new_releases/event"
	. "Spotify_new_releases/database"
	. "Spotify_new_releases/cookie"
)

var (
	errURI = "https://newreleases.tk" //os.Getenv("ERR_URI")
	sucURI = "https://newreleases.tk/user" //os.Getenv("SUC_URI")
)


func RedirectHandler(w http.ResponseWriter, r *http.Request, mydbmap *MyDbMap) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers","Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	
	// use the same state string here that you used to generate the URL
	fmt.Println("redirect")
	fmt.Println(sucURI)
	fmt.Println("redirect")
	// create client and get token
	client, token, r, err := CreateMyClient(r)
	
	// get user id from client
	userId, err := client.GetCurrentUserId()
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, errURI, 302)
	}

	// check if the client is already in database.
	// if not, then create a playlist id and insert user info into database
	var playlistId string
	ifExists, user, err := mydbmap.ExistUser(userId)
	if !ifExists {
		playlist, err := client.CreatePlaylistForUser(userId)
		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, errURI, 302)
		}
		if err := client.SetConfig(playlist.ID); err != nil {
			fmt.Println(err)
			http.Redirect(w, r, errURI, 302)
		}
		playlistId = string(playlist.ID)
		mydbmap.InsertUser(userId, playlistId, token)
		if err := GetFollowingArtistsAndInsertRelations(mydbmap, userId, token); err != nil {
			fmt.Println(err)
			http.Redirect(w, r, errURI, 302)
		}
	} else {
		playlistId = user.PlaylistId
	}

	// set cookies
	w, err = SetCookie(w, token)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, errURI, 302)
	}
	fmt.Println(w)
	// http.Redirect(w, r, "http://localhost:8080/user/"+token.RefreshToken, 301)
	http.Redirect(w, r, sucURI+"/"+token.AccessToken, 301)
}
