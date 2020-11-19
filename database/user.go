package database

import (
	//"database/sql"
	_ "github.com/go-sql-driver/mysql"
	//"github.com/go-gorp/gorp"
	"golang.org/x/oauth2"
	"fmt"
	"time"
)

type UserInfo struct {
	UserId			string 		`db:userId`
	AccessToken		string 		`db:accessToken`
	TokenType		string 		`db:tokenType`
	RefreshToken	string 		`db:refreshToken`
	Expiry			time.Time 	`db:expiry`
	PlaylistId		string 		`db:playlistId`
}

func (d *MyDbMap) ExistUser(userId string) (bool, UserInfo, error) {
	var user UserInfo
	err := dbmap.SelectOne(&user, "select * from User where userId=?", userId)
	if err != nil {
		fmt.Println(err)
		return false, user, err
	}
	return true, user, nil
}

func (d *MyDbMap) InsertUser(userId, playlistId string, Token *oauth2.Token) error {
	count, err := d.DbMap.SelectInt("select count(*) from User where userId = ?", userId)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if count == 0 {
		//err = dbmap.Insert(&UserInfo{UserId: userId, AccessToken: Token.AccessToken, TokenType: Token.TokenType, RefreshToken: Token.RefreshToken, Expiry: Token.Expiry})
		err = d.DbMap.Insert(&UserInfo{userId, Token.AccessToken, Token.TokenType, Token.RefreshToken, Token.Expiry, playlistId})
	}
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (d *MyDbMap) GetAllUsers() ([]UserInfo, error) {
	var users []UserInfo
	if _, err := d.DbMap.Select(&users, "select * from User"); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return users, nil
}

func (d *MyDbMap) UpdateUser(userId, playlistId string, Token *oauth2.Token) error {
	user := UserInfo{userId, Token.AccessToken, Token.TokenType, Token.RefreshToken, Token.Expiry, playlistId}
	if _, err := d.DbMap.Update(&user); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
