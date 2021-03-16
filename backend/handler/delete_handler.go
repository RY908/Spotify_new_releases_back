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
	Status 	int 			`json:"status"`
	Result 	string 			`json:"result"`
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
	// TODO: status 400
	exists, user, err := GetUser(r, mydbmap, token)
	if err != nil {
		// TODO: status 500
		response := DeleteResponse{400, "failed", []ArtistInfo{}}
		res, err := json.Marshal(response)
		fmt.Println(err)
		w.Write(res)
	}

	// if the user is not in database then return response without artist information
	// if the user is in database then delete the artists the user requests and return artists
	if exists == false {
		// Todo: status 401
		response := DeleteResponse{200, "redirect", []ArtistInfo{}}
		res, err := json.Marshal(response)
		if err != nil {
			fmt.Println(err)
		}
		w.Write(res)
	} else {
		var request DeleteRequest
		json.NewDecoder(r.Body).Decode(&request)
		artistIds := request.ArtistIds

		if err := mydbmap.DeleteRelationFromRequest(user.UserId, artistIds); err != nil {
			// TODO: w.write
			fmt.Println(err)
		}

		artists, err := mydbmap.GetArtistsFromUserId(user.UserId)
		if err != nil {
			// TODO: w.write
			fmt.Println(err)
		}
		response := DeleteResponse{200, "success", artists}
		res, err := json.Marshal(response)
		if err != nil {
			// TODO: w.write
			fmt.Println(err)
		}
		w.Write(res)
	}
}
