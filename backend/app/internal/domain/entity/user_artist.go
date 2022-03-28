package entity

import "github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/schema"

type UserArtist struct {
	ID          string
	Name        string
	Url         string
	IconUrl     string
	IsFollowing bool
}

func NewUserArtist(userArtist *schema.UserArtist) *UserArtist {
	return &UserArtist{
		ID:          userArtist.ID,
		Name:        userArtist.Name,
		Url:         userArtist.Url,
		IconUrl:     userArtist.IconUrl,
		IsFollowing: userArtist.IsFollowing,
	}
}

func NewUserArtists(userArtists []*schema.UserArtist) []*UserArtist {
	var userArtistsEntity []*UserArtist
	for _, userArtist := range userArtists {
		userArtistsEntity = append(userArtistsEntity, NewUserArtist(userArtist))
	}
	return userArtistsEntity
}
