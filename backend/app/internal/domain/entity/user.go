package entity

import (
	"golang.org/x/oauth2"
	"time"

	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/schema"
)

type User struct {
	ID            string
	AccessToken   string
	TokenType     string
	RefreshToken  string
	Expiry        time.Time
	PlaylistID    string
	IfRemixAdd    bool
	IfAcousticAdd bool
}

func NewUser(user *schema.User) *User {
	return &User{
		ID:            user.Id,
		AccessToken:   user.AccessToken,
		TokenType:     user.TokenType,
		RefreshToken:  user.RefreshToken,
		Expiry:        user.Expiry,
		PlaylistID:    user.PlaylistId,
		IfRemixAdd:    user.IfRemixAdd,
		IfAcousticAdd: user.IfAcousticAdd,
	}
}

func (u *User) UpdateUserByToken(token *oauth2.Token) *User {
	return &User{
		ID:            u.ID,
		AccessToken:   token.AccessToken,
		TokenType:     token.TokenType,
		RefreshToken:  token.RefreshToken,
		Expiry:        token.Expiry,
		PlaylistID:    u.PlaylistID,
		IfRemixAdd:    u.IfRemixAdd,
		IfAcousticAdd: u.IfAcousticAdd,
	}
}
