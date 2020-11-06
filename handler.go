package main

import (
	"fmt"
	"net/http"
	//"log"
	//"os"
	"time"
	. "Spotify_new_releases/spotify"
	. "Spotify_new_releases/session"
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

	http.Redirect(w, r, "/home", 301)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("home")

	userId, token, playlistId := GetTokenFromSession(r)
	fmt.Println(token)
	
	client := CreateMyClientFromToken(token)
	fmt.Println(client)
	
	artists, newToken := client.GetRecentlyPlayedArtists()
	fmt.Println(artists)

	timestamp := time.Now()
	fmt.Println(newToken)

	for _, artist := range artists {
		//fmt.Println(artist)
		name := artist.SimpleArtist.Name
		artistId := artist.SimpleArtist.ID.String()
		url := artist.SimpleArtist.ExternalURLs["spotify"]
		iconUrl := artist.Images[0].URL
		
		fmt.Println(artist)
		err := mydbmap.InsertArtist(artistId, name, url, iconUrl)
		if err != nil {
			fmt.Println(62, name)
			fmt.Println(err)
		}

		err = mydbmap.InsertRelation(artistId, userId, timestamp)
		if err != nil {
			fmt.Println(68, name)
			fmt.Println(err)
		}
	}
	err := mydbmap.UpdateUser(userId, playlistId, newToken)
	if err != nil {
		fmt.Println(err)
	}

	newReleases, err := client.GetNewReleases(mydbmap, userId)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(newReleases)
}
