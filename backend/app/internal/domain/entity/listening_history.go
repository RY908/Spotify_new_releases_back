package entity

import (
	"time"

	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/schema"
)

type ListeningHistory struct {
	Id          int64
	ArtistId    string
	UserId      string
	ListenCount int64
	Timestamp   time.Time
	IsFollowing bool
}

func NewListeningHistory(listeningHistory *schema.ListeningHistory) ListeningHistory {
	return ListeningHistory{
		ArtistId:    listeningHistory.ArtistId,
		UserId:      listeningHistory.UserId,
		ListenCount: listeningHistory.ListenCount,
		Timestamp:   listeningHistory.Timestamp,
		IsFollowing: listeningHistory.IsFollowing,
	}
}
