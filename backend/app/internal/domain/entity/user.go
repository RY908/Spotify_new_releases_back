package entity

import (
	"time"

	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/schema"
)

type User struct {
	Id            string
	AccessToken   string
	TokenType     string
	RefreshToken  string
	Expiry        time.Time
	PlaylistId    string
	IfRemixAdd    bool
	IfAcousticAdd bool
}

func NewUser(user *schema.User) User {
	return User{
		Id:            user.Id,
		AccessToken:   user.AccessToken,
		TokenType:     user.TokenType,
		RefreshToken:  user.RefreshToken,
		Expiry:        user.Expiry,
		PlaylistId:    user.PlaylistId,
		IfRemixAdd:    user.IfRemixAdd,
		IfAcousticAdd: user.IfAcousticAdd,
	}
}
