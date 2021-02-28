package handler

import (
	"fmt"
	"net/http"
	"encoding/json"
	//. "Spotify_new_releases/spotify"
	//. "Spotify_new_releases/session"
	. "Spotify_new_releases/database"
	. "Spotify_new_releases/cookie"
)

type UserResponse struct {
	Status 	int 			`json:"status"`
	Result 	string 			`json:"result"`
	Artists []ArtistInfo 	`json:"artists"`
}

func UserHandler(w http.ResponseWriter, r *http.Request, mydbmap *MyDbMap) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers","Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// check if the user is in database.
	// if not, then set status as redirect
	// if existed, then get all the artists and send them to frontend.
	exists, user, err := GetUser(r, mydbmap)
	if err != nil {
		response := UserResponse{400, "failed", []ArtistInfo{}}
		res, err := json.Marshal(response)
		fmt.Println(err)
		w.Write(res)
	}
	if exists == false {
		response := UserResponse{200, "redirect", []ArtistInfo{}}
		res, err := json.Marshal(response)
		if err != nil {
			fmt.Println(err)
		}
		w.Write(res)
	} else {
		artists, err := mydbmap.GetArtistsFromUserId(user.UserId)
		if err != nil {
			fmt.Println(err)
		}
		response := UserResponse{200, "success", artists}
		res, err := json.Marshal(response)
		if err != nil {
			fmt.Println(err)
		}
		w.Write(res)
	}
}	

