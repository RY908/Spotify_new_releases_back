package database

import (
	"fmt"
	//"database/sql"
	_ "github.com/go-sql-driver/mysql"
	//"github.com/go-gorp/gorp"
	"golang.org/x/oauth2"
	"time"
)

type UserInfo struct {
	UserId        string    `db:userId`
	AccessToken   string    `db:accessToken`
	TokenType     string    `db:tokenType`
	RefreshToken  string    `db:refreshToken`
	Expiry        time.Time `db:expiry`
	PlaylistId    string    `db:playlistId`
	IfRemixAdd    bool      `db:ifRemixAdd`
	IfAcousticAdd bool      `db:ifAcousticAdd`
}

// check if the user is in database
func (d *MyDbMap) ExistUser(userId string) (bool, UserInfo, error) {
	var user UserInfo
	err := d.DbMap.SelectOne(&user, "select * from User where userId=?", userId)
	if err != nil {
		return false, user, fmt.Errorf("select user: %w", ErrUserNotFound)
	}
	return true, user, nil
}

// insert user in database
func (d *MyDbMap) InsertUser(userId, playlistId string, Token *oauth2.Token) error {
	count, err := d.DbMap.SelectInt("select count(*) from User where userId = ?", userId)
	if err != nil {
		return err
	}
	if count == 0 {
		//err = dbmap.Insert(&UserInfo{UserId: userId, AccessToken: Token.AccessToken, TokenType: Token.TokenType, RefreshToken: Token.RefreshToken, Expiry: Token.Expiry})
		err = d.DbMap.Insert(&UserInfo{userId, Token.AccessToken, Token.TokenType, Token.RefreshToken, Token.Expiry, playlistId, true, true})
	}
	if err != nil {
		return fmt.Errorf("insert user: %w", ErrUserAlreadyExisted)
	}
	return nil
}

// get all users in database
func (d *MyDbMap) GetAllUsers() ([]UserInfo, error) {
	var users []UserInfo
	if _, err := d.DbMap.Select(&users, "select * from User"); err != nil {
		return nil, fmt.Errorf("select all users: %w", ErrUserNotFound)
	}

	return users, nil
}

// update user's auth information
func (d *MyDbMap) UpdateUser(userId, playlistId string, Token *oauth2.Token) error {
	/*user := UserInfo{userId, Token.AccessToken, Token.TokenType, Token.RefreshToken, Token.Expiry, playlistId, true, true}
	if _, err := d.DbMap.Update(&user); err != nil {
		return err
	}*/
	if _, err := d.DbMap.Exec("update User set accessToken = ?, tokenType = ?, refreshToken = ?, expiry = ? where userId = ?", Token.AccessToken, Token.TokenType, Token.RefreshToken, Token.Expiry, userId); err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

func (d *MyDbMap) UpdateIfAdd(userId string, ifRemixAdd, ifAcousticAdd bool) error {
	if _, err := d.DbMap.Exec("update User set ifRemixAdd = ?, ifAcousticAdd = ? where userId = ?", ifRemixAdd, ifAcousticAdd, userId); err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil

}
