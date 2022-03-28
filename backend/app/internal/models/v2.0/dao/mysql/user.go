package mysql

import (
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/schema"
	"github.com/go-gorp/gorp"
)

type user struct {
	db *gorp.DbMap
}

func (u *user) InsertUser(user *schema.User) error {
	if err := u.db.Insert(user); err != nil {
		return err
	}

	return nil
}

func (u *user) GetUser(userID string) (*schema.User, error) {
	var user schema.User
	if err := u.db.SelectOne(&user, "select * from User where userId=?", userID); err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *user) GetAllUsers() ([]*schema.User, error) {
	var users []*schema.User
	if _, err := u.db.Select(&users, "select * from User"); err != nil {
		return nil, err
	}
	return users, nil
}

func (u *user) UpdateUser(user *schema.User) error {
	if _, err := u.db.Update(user); err != nil {
		return err
	}
	return nil
}
