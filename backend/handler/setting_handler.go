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
	Status 			int 	`json:"status"`
	Result 			string 	`json:"result"`
	IfRemixAdd 		bool 	`json:"ifRemixAdd"`
	IfAcousticAdd 	bool 	`json:"ifAcousticAdd"`
}

type SettingEditResponse struct {
	Status 			int 	`json:"status"`
	Result 			string 	`json:"result"`
}

func SettingHandler(w http.ResponseWriter, r *http.Request, mydbmap *MyDbMap) {
	w.Header().Set("Access-Control-Allow-Origin", accessControlAllowOrigin)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers","Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	fmt.Println("setting")

	exists, user, err := GetUser(r, mydbmap)
	if err != nil {
		response := SettingResponse{400, "failed", false, false}
		res, err := json.Marshal(response)
		fmt.Println(err)
		w.Write(res)
	}

	if exists == false {
		response := SettingResponse{200, "redirect", false, false}
		res, err := json.Marshal(response)
		if err != nil {
			fmt.Println(err)
		}
		w.Write(res)
	} else {
		response := SettingResponse{200, "success", user.IfRemixAdd, user.IfAcousticAdd}
		res, err := json.Marshal(response)
		if err != nil {
			fmt.Println(err)
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

	exists, user, err := GetUser(r, mydbmap)
	if err != nil {
		response := SettingEditResponse{500, "failed"}
		res, err := json.Marshal(response)
		fmt.Println(err)
		w.Write(res)
	}
	if exists == false {
		response := SettingEditResponse{200, "redirect"}
		res, err := json.Marshal(response)
		if err != nil {
			fmt.Println(err)
		}
		w.Write(res)
	} else {
		var request SettingEditRequest
		var response SettingEditResponse
		json.NewDecoder(r.Body).Decode(&request)
		ifRemixAdd := request.IfRemixAdd
		ifAcousticAdd := request.IfAcousticAdd

		if err := mydbmap.UpdateIfAdd(user.UserId, ifRemixAdd, ifAcousticAdd); err != nil {
			fmt.Println(err)
			response = SettingEditResponse{500, "failed"}
		} else {
			response = SettingEditResponse{200, "success"}
		}
		
		res, err := json.Marshal(response)
		if err != nil {
			fmt.Println(err)
		}
		w.Write(res)
	}
}
