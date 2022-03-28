package dao

import (
	"errors"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/schema"
	"github.com/go-gorp/gorp"
	"time"
)

var factoryFunc func(db *gorp.DbMap) Factory

func RegisterFactory(f func(db *gorp.DbMap) Factory) {
	factoryFunc = f
}

func NewDBManager(db *gorp.DbMap) (Factory, error) {
	if factoryFunc == nil {
		return nil, errors.New("no factory found, register it by importing implementation of dao.Factory")
	}

	return factoryFunc(db), nil
}

type Factory interface {
	ArtistDAO() Artist
	ListeningHistoryDAO() ListeningHistory
	UserDAO() User
	UserArtistsDAO() UserArtists
}

type Artist interface {
	InsertArtist(artist *schema.Artist) error
	InsertArtists(artists []*schema.Artist) error
}

type ListeningHistory interface {
	GetAllListeningHistory() ([]*schema.ListeningHistory, error)
	GetListeningHistory(artistId, userId string) (*schema.ListeningHistory, error)
	InsertHistory(history *schema.ListeningHistory) error
	UpdateHistory(artistId, userId string, count int64, isFollowing bool, timestamp time.Time) error
	DeleteHistoryByTimestamp(userId string, timestamp time.Time) error
	DeleteHistoryByArtistId(userId, artistId string) error
}

type User interface {
	InsertUser(user *schema.User) error
	GetUser(userId string) (*schema.User, error)
	GetAllUsers() ([]*schema.User, error)
	UpdateUser(user *schema.User) error
}

type UserArtists interface {
	GetArtistsByUserID(userID string) ([]*schema.UserArtist, error)
}
