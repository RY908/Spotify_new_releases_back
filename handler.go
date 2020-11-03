package main

import (
	"fmt"
	"net/http"
	"log"
	//"os"
	"time"
)


func redirectHandler(w http.ResponseWriter, r *http.Request) {
	// use the same state string here that you used to generate the URL
	fmt.Println("/handle")
	token, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusNotFound)
		return
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	client := auth.NewClient(token)
	user, err := client.CurrentUser()
	if err != nil {
		fmt.Println(err)
	}

	
	session, _ := store.Get(r, session_name)
	session.Values["user"] = UserSession{ID: user.ID, Token:*token}
	err = session.Save(r, w)

	insertUser(user.ID, token)

	http.Redirect(w, r, "/home", 301)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("home")

	userId, token := getTokenFromSession(r)
	fmt.Println(token)
	client := auth.NewClient(&token)
	fmt.Println(client)
	artists, newToken := getRecentlyPlayedArtists(client)

	timestamp := time.Now()
	fmt.Println(newToken)

	for _, artist := range artists {
		//fmt.Println(artist)
		name := artist.SimpleArtist.Name
		artistId := artist.SimpleArtist.ID.String()
		url := artist.SimpleArtist.ExternalURLs["spotify"]
		iconUrl := artist.Images[0].URL
		
		err := insertArtist(artistId, name, url, iconUrl)
		if err != nil {
			fmt.Println(62, name)
			fmt.Println(err)
		}

		err = insertRelation(artistId, userId, timestamp)
		if err != nil {
			fmt.Println(68, name)
			fmt.Println(err)
		}
	}
	err := updateUser(userId, newToken)
	if err != nil {
		fmt.Println(err)
	}

	newReleases, err := getNewReleases(client, userId)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(newReleases)
}
