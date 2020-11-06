package main

import (
	"fmt"
	"net/http"
	//"log"
	//"os"
	//"time"
	. "Spotify_new_releases/spotify"
	. "Spotify_new_releases/session"
	. "Spotify_new_releases/event"
	//. "Spotify_new_releases/database"
)


func redirectHandler(w http.ResponseWriter, r *http.Request) {
	// use the same state string here that you used to generate the URL
	fmt.Println("/handle")
	/*
	token, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusNotFound)
		return
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	client := auth.NewClient(token)*/
	client, token, r := CreateMyClient(r)
	//client = client.Client
	userId, err := client.GetCurrentUserId()
	if err != nil {
		fmt.Println(err)
	}

	playlist, err := client.CreatePlaylistForUser(userId)
	if err != nil {
		fmt.Println(err)
	}
	playlistId := string(playlist.ID)
	
	session, _ := Store.Get(r, Session_name)
	session.Values["user"] = UserSession{ID: userId, Token:*token, PlaylistId: playlistId}
	err = session.Save(r, w)

	mydbmap.InsertUser(userId, playlistId, token)
	if err := GetFollowingArtistsAndInsertRelations(mydbmap, userId, token); err != nil {
		fmt.Println(err)
	}

	http.Redirect(w, r, "/home", 301)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("home")
	if err := UpdateRelation(mydbmap); err != nil {
		fmt.Println(err)
	}
	if err := UpdatePlaylist(mydbmap); err != nil {
		fmt.Println(err)
	}
	
}
