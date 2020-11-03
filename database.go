package main

import (
	//"database/sql"
	_ "github.com/go-sql-driver/mysql"
	//"github.com/go-gorp/gorp"
	"golang.org/x/oauth2"
	"fmt"
	"time"
)

type ArtistInfo struct {
    ArtistId	string 		`db:"artistId"`
	Name 		string 		`db:"name"`
	Url 		string 		`db:"url"`
	IconUrl		string 		`db:iconUrl`
}

type ListenTo struct {
	ListenId	int64 	`db:listenId`
	ArtistId	string 	`db:artistId`
	UserId		string 	`db:userId`
	Timestamp	time.Time	`db:"timestamp"`
}

type UserInfo struct {
	UserId			string 		`db:userId`
	AccessToken		string 		`db:accessToken`
	TokenType		string 		`db:tokenType`
	RefreshToken	string 		`db:refreshToken`
	Expiry			time.Time 	`db:expiry`
}

func insertArtist(artistId, name, url, iconUrl string) error {
	count, err := dbmap.SelectInt("select count(*) from Artist where artistId = ?", artistId) // check if artist already exists in database
	if err != nil {
		fmt.Println(err)
		return err
	}

	// if artist does not exist, then insert artist into database
	if count == 0 {
		err = dbmap.Insert(&ArtistInfo{ArtistId:artistId, Name:name, Url:url, IconUrl:iconUrl})
	}
	
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func insertRelation(artistId, userId string, timestamp time.Time) error {
	count, err := dbmap.SelectInt("select count(*) from ListenTo where artistId = ? and userId = ?", artistId, userId)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if count == 0 {
		err = dbmap.Insert(&ListenTo{ArtistId:artistId, UserId: userId, Timestamp:timestamp})
	}
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func insertUser(userId string, Token *oauth2.Token) error {
	count, err := dbmap.SelectInt("select count(*) from User where userId = ?", userId)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if count == 0 {
		//err = dbmap.Insert(&UserInfo{UserId: userId, AccessToken: Token.AccessToken, TokenType: Token.TokenType, RefreshToken: Token.RefreshToken, Expiry: Token.Expiry})
		err = dbmap.Insert(&UserInfo{userId, Token.AccessToken, Token.TokenType, Token.RefreshToken, Token.Expiry})
	}
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func deleteRelation(userId string, timestamp time.Time) error {
	res, err := dbmap.Exec("delete from ListenTo where userId = ? and timestamp = ?", userId, timestamp)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(res)
	return nil
}

func getAllUsers() ([]UserInfo, error) {
	var users []UserInfo
	_, err := dbmap.Select(&users, "select * from users")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return users, nil
}

func updateUser(userId string, Token *oauth2.Token) error {
	user := UserInfo{userId, Token.AccessToken, Token.TokenType, Token.RefreshToken, Token.Expiry}
	_, err := dbmap.Update(&user)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func getArtistsFromUserId(userId string) ([]ArtistInfo, error) {
	var artists []ArtistInfo
	cmd := "select Artist.artistId, Artist.name, Artist.url, Artist.iconUrl from Artist inner join ListenTo on Artist.artistId = ListenTo.artistId where ListenTo.userId = ?"
	_, err := dbmap.Select(&artists, cmd, userId)
	if err != nil {
		fmt.Println(err)
		return nil, err  
	}
	return artists, nil
}