package schema

import "time"

type ListeningHistory struct {
	ID          int64     `db:"listenId"`
	ArtistID    string    `db:"artistId"`
	UserID      string    `db:"userId"`
	ListenCount int64     `db:"listenCount"`
	Timestamp   time.Time `db:"timestamp"`
	IsFollowing bool      `db:"isFollowing"`
}
