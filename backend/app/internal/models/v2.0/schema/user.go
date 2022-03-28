package schema

import "time"

type User struct {
	ID            string    `db:"userId"`
	AccessToken   string    `db:"accessToken"`
	TokenType     string    `db:"tokenType"`
	RefreshToken  string    `db:"refreshToken"`
	Expiry        time.Time `db:"expiry"`
	PlaylistId    string    `db:"playlistId"`
	IfRemixAdd    bool      `db:"ifRemixAdd"`
	IfAcousticAdd bool      `db:"ifAcousticAdd"`
}
