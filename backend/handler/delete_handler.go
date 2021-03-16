package handler

import (
	"fmt"
	"net/http"
	"encoding/json"
	. "Spotify_new_releases/database"
	. "Spotify_new_releases/cookie"
)

type DeleteRequest struct {
	ArtistIds []string	`json:"artistsId"`
}

type DeleteResponse struct {
	Artists []ArtistInfo 	`json:"artists"`
}

func DeleteHandler(w http.ResponseWriter, r *http.Request, mydbmap *MyDbMap) {
	w.Header().Set("Access-Control-Allow-Origin", accessControlAllowOrigin)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers","Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	fmt.Println("delete")

	// get user from cookie
	token, err := GetToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
        return
	}
	exists, user, err := GetUser(r, mydbmap, token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
        return
	}

	// if the user is not in database then return response without artist information
	// if the user is in database then delete the artists the user requests and return artists
	if exists == false {
		http.Error(w, err.Error(), http.StatusUnauthorized)
        return
	} else {
		var request DeleteRequest
		json.NewDecoder(r.Body).Decode(&request)
		artistIds := request.ArtistIds

		if err := mydbmap.DeleteRelationFromRequest(user.UserId, artistIds); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		artists, err := mydbmap.GetArtistsFromUserId(user.UserId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response := DeleteResponse{artists}
		res, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(res)
	}
}
