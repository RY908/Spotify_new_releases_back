package schema

import "time"

type ListeningHistory struct {
	Id          int64     `db:"listenId"`
	ArtistId    string    `db:"artistId"`
	UserId      string    `db:"userId"`
	ListenCount int64     `db:"listenCount"`
	Timestamp   time.Time `db:"timestamp"`
	IsFollowing bool      `db:"isFollowing"`
}
