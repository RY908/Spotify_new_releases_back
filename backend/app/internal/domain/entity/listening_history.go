package entity

import (
	"time"

	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/schema"
)

type ListeningHistory struct {
	ID          int64
	ArtistID    string
	UserID      string
	ListenCount int64
	Timestamp   time.Time
	IsFollowing bool
}

func NewListeningHistory(listeningHistory *schema.ListeningHistory) *ListeningHistory {
	return &ListeningHistory{
		ID:          listeningHistory.ID,
		ArtistID:    listeningHistory.ArtistID,
		UserID:      listeningHistory.UserID,
		ListenCount: listeningHistory.ListenCount,
		Timestamp:   listeningHistory.Timestamp,
		IsFollowing: listeningHistory.IsFollowing,
	}
}
