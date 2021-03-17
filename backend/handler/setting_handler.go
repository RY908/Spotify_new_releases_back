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

type SettingEditRequest struct {
	IfRemixAdd 		bool `json:"ifRemixAdd"`
	IfAcousticAdd 	bool `json:"ifAcousticAdd"`
}

type SettingResponse struct {
	IfRemixAdd 		bool 	`json:"ifRemixAdd"`
	IfAcousticAdd 	bool 	`json:"ifAcousticAdd"`
}

func SettingHandler(w http.ResponseWriter, r *http.Request, mydbmap *MyDbMap) {
	w.Header().Set("Access-Control-Allow-Origin", accessControlAllowOrigin)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers","Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	fmt.Println("setting")

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

	if exists == false {
		http.Error(w, err.Error(), http.StatusUnauthorized)
        return
	} else {
		response := SettingResponse{user.IfRemixAdd, user.IfAcousticAdd}
		res, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(res)
	}
}

func SettingEditHandler(w http.ResponseWriter, r *http.Request, mydbmap *MyDbMap) {
	w.Header().Set("Access-Control-Allow-Origin", accessControlAllowOrigin)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers","Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	
	fmt.Println("change setting")

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

	if exists == false {
		http.Error(w, err.Error(), http.StatusUnauthorized)
        return
	} else {
		var request SettingEditRequest
		// var response SettingEditResponse
		json.NewDecoder(r.Body).Decode(&request)
		ifRemixAdd := request.IfRemixAdd
		ifAcousticAdd := request.IfAcousticAdd

		if err := mydbmap.UpdateIfAdd(user.UserId, ifRemixAdd, ifAcousticAdd); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
