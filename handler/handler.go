package handler

import (
	"fmt"
	"net/http"
	//"log"
	//"os"
	"golang.org/x/oauth2"
	//"io/ioutil"
	"time"
	"encoding/json"
	//"html/template"
	. "Spotify_new_releases/spotify"
	//. "Spotify_new_releases/session"
	. "Spotify_new_releases/event"
	. "Spotify_new_releases/database"
)

var layout = "2006-01-02 15:04:05"

type Request struct {
	Token string `json:token`
}

type Response struct {
	Status 	int 			`json:"status"`
	Result 	string 			`json:"result"`
	Artists []ArtistInfo 	`json:"artists"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("login")
	fmt.Println(r.Header["Referer"])
	url := GetURL()
	fmt.Println(url)
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Headers","Content-Type")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
	/*t := template.Must(template.ParseFiles("templates/index.html"))
	if err := t.Execute(w, url); err != nil {
		fmt.Println(err)
	}*/
	/*response := Response{200, "ok"}
        res, err := json.Marshal(response)
        if err != nil {
                fmt.Println(err)
        }
	w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("Access-Control-Allow-Headers","Content-Type")
	w.Write(res)
	fmt.Println(res)*/
	http.Redirect(w, r, url, 302)
}


func RedirectHandler(w http.ResponseWriter, r *http.Request, mydbmap *MyDbMap) {
	// use the same state string here that you used to generate the URL
	fmt.Println("/handle")

	// create client and get token
	client, token, r, err := CreateMyClient(r)
	
	// get user id from client
	userId, err := client.GetCurrentUserId()
	if err != nil {
		fmt.Println(err)
	}

	// check if the client is already in database.
	// if not, then create a playlist id and insert user info into database
	var playlistId string
	ifExists, user, err := mydbmap.ExistUser(userId)
	if !ifExists {
		playlist, err := client.CreatePlaylistForUser(userId)
		if err != nil {
			fmt.Println(err)
		}
		playlistId = string(playlist.ID)
		mydbmap.InsertUser(userId, playlistId, token)
		if err := GetFollowingArtistsAndInsertRelations(mydbmap, userId, token); err != nil {
			fmt.Println(err)
		}
	} else {
		playlistId = user.PlaylistId
	}

	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers","Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// set cookies
    http.SetCookie(w, &http.Cookie{
        Name: "token",
        Value: token.AccessToken,
     	Path: "/",
	})
	http.SetCookie(w, &http.Cookie{
		Name: "tokenType",
		Value: token.TokenType,
		Path: "/",
	})
	http.SetCookie(w, &http.Cookie{
		Name: "refreshToken",
		Value: token.RefreshToken,
		Path: "/",
	})
	http.SetCookie(w, &http.Cookie{
		Name: "expiry",
		Value: (token.Expiry).Format(layout),
		Path: "/",
	})

	http.Redirect(w, r, "http://localhost:8080/callback/"+token.RefreshToken, 302)
}

func UserHandler(w http.ResponseWriter, r *http.Request, mydbmap *MyDbMap) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers","Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// get cookies
	tokenCookie, err := r.Cookie("token")
	if err != nil {
		fmt.Println(err)
	}
	refreshTokenCookie, err := r.Cookie("refreshToken")
	if err != nil {
		fmt.Println(err)
	}
	tokenTypeCookie, err := r.Cookie("tokenType")
    if err != nil {
		fmt.Println(err)
	}
	expiryCookie, err := r.Cookie("expiry")
	if err != nil {
		fmt.Println(err)
	}
	accessToken := tokenCookie.Value
	refreshToken := refreshTokenCookie.Value
	tokenType := tokenTypeCookie.Value
	expiryString := expiryCookie.Value
	expiry, _ := time.Parse(layout, expiryString)

	// get token, client and user id
	token := oauth2.Token{AccessToken:accessToken, TokenType:tokenType, RefreshToken:refreshToken, Expiry:expiry}
	client := CreateMyClientFromToken(token)
	userId, err := client.GetCurrentUserId()
	if err != nil {
		fmt.Println(err)
	}

	// check if the user is in database.
	// if not, then set status as redirect
	// if existed, then get all the artists and send them to frontend.
	exists, user, err := mydbmap.ExistUser(userId)
	if err != nil {
		fmt.Println(err)
	}
	if exists == false {
		response := Response{200, "redirect", []ArtistInfo{}}
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
		response := Response{200, "success", artists}
		res, err := json.Marshal(response)
		if err != nil {
			fmt.Println(err)
		}
		w.Write(res)
	}
}	
