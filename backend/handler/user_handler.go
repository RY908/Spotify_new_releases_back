package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	//. "Spotify_new_releases/spotify"
	//. "Spotify_new_releases/session"
	. "Spotify_new_releases/backend/cookie"
	. "Spotify_new_releases/backend/database"
)

type UserResponse struct {
	Artists []ArtistRes `json:"artists"`
}

func UserHandler(w http.ResponseWriter, r *http.Request, mydbmap *MyDbMap) {
	w.Header().Set("Access-Control-Allow-Origin", accessControlAllowOrigin)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	fmt.Println(r)
	// check if the user is in database.
	// if not, then set status as redirect
	// if existed, then get all the artists and send them to frontend.
	token, err := GetToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	exists, user, err := GetUser(r, mydbmap, token)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if exists == false {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	} else {
		artists, err := mydbmap.GetArtistsFromUserId(user.UserId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response := UserResponse{artists}
		res, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(res)
	}
}
