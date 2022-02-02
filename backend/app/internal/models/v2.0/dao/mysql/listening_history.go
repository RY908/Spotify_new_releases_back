package mysql

import (
	"fmt"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/schema"
	"github.com/go-gorp/gorp"
	"time"
)

type listeningHistory struct {
	db *gorp.DbMap
}

func (l *listeningHistory) GetAllListeningHistory() ([]*schema.ListeningHistory, error) {
	var history []*schema.ListeningHistory
	if _, err := l.db.Select(&history, "select * from ListenTo"); err != nil {
		return nil, fmt.Errorf("get all relations: %w", err)
	}
	return history, nil
}

func (l *listeningHistory) GetListeningHistory(artistId, userId string) (*schema.ListeningHistory, error) {
	var history *schema.ListeningHistory
	if err := l.db.SelectOne(&history, "select * from ListenTo where artistId = ? and userID = ?", artistId, userId); err != nil {
		return nil, err
	}
	return history, nil
}

func (l *listeningHistory) InsertHistory(history *schema.ListeningHistory) error {
	if err := l.db.Insert(history); err != nil {
		return err
	}
	return nil
}

func (l *listeningHistory) UpdateHistory(artistId, userId string, count int64, isFollowing bool, timestamp time.Time) error {
	if _, err := l.db.Exec("update ListenTo set listenCount = listenCount+?, timestamp = ?, ifFollowing = ? where artistId = ? and userId = ?", count, timestamp, isFollowing, artistId, userId); err != nil {
		return err
	}
	return nil
}

func (l *listeningHistory) DeleteHistoryByTimestamp(userId string, timestamp time.Time) error {
	if _, err := l.db.Exec("delete from ListenTo where userId = ? and timestamp < ? and ifFollowing = false", userId, timestamp); err != nil {
		return err
	}
	return nil
}

func (l *listeningHistory) DeleteHistoryByArtistId(userId, artistId string) error {
	if _, err := l.db.Exec("delete from ListenTo where userId = ? and artistId = ?", userId, artistId); err != nil {
		return err
	}
	return nil
}
