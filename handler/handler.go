package handler

import (
	"fmt"
	"net/http"
	//"log"
	//"os"
	//"golang.org/x/oauth2"
	//"io/ioutil"
	//"time"
	"encoding/json"
	//"html/template"
	. "Spotify_new_releases/spotify"
	. "Spotify_new_releases/session"
	. "Spotify_new_releases/event"
	. "Spotify_new_releases/database"
)

type Request struct {
	Token string `json:token`
}

type Response struct {
	Status int 		`json:status`
	Result string 	`json:result`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("login")
	fmt.Println(r.Header["Referer"])
	url := GetURL()
	fmt.Println(url)
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
	http.Redirect(w, r, url, 301)
}


func RedirectHandler(w http.ResponseWriter, r *http.Request, mydbmap *MyDbMap) {
	// use the same state string here that you used to generate the URL
	fmt.Println("/handle")

	/*body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	bytes := []byte(body)
	var request Request
	json.Unmarshal(bytes, &request)
	
	token := request.Token
	fmt.Println(token)
	*/
	// token = oauth2.Token(token)

	// create client and get token
	client, token, r, err := CreateMyClient(r)
	// client, token, r, err := CreateMyClientFromCode(r)
	//client := CreateMyClientFromToken(*token)
	
	userId, err := client.GetCurrentUserId()
	if err != nil {
		fmt.Println(err)
	}

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
	
	session, _ := Store.Get(r, Session_name)
	session.Values["user"] = UserSession{ID: userId, Token:*token, PlaylistId: playlistId}
	err = session.Save(r, w)
	
	// http.Redirect(w, r, "/home", 301)
	response := Response{200, "ok"}
	res, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers","Content-Type")
	fmt.Println(res)
	cookie := &http.Cookie{
          Name: "hoge",
          Value: "bar",
     	}
     	http.SetCookie(w, cookie)
	http.Redirect(w, r, "http://localhost:8080/callback"+userId, 301)
	w.Write(res)
}

func HomeHandler(w http.ResponseWriter, r *http.Request, mydbmap *MyDbMap) {
	fmt.Println("home")
	fmt.Println(r)

	response := Response{200, "ok"}
	res, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
	/*t := template.Must(template.ParseFiles("templates/home.html"))
	if err := t.Execute(w, time.Now()); err != nil {
		fmt.Println(err)
	}
	if err := UpdateRelation(mydbmap); err != nil {
		fmt.Println(err)
	}*/
	/*if err := UpdatePlaylist(mydbmap); err != nil {
		fmt.Println(err)
	}*/
	
}
