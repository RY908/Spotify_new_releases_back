package repository

import (
	"database/sql"
	"errors"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/entity"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/dao"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/schema"
	"time"
)

func NewListeningHistoryRepository() *ListeningHistoryRepository {
	return &ListeningHistoryRepository{}
}

type ListeningHistoryRepository struct{}

func (r *ListeningHistoryRepository) InsertOrUpdateListeningHistory(factory dao.Factory, history entity.ListeningHistory, counter map[string]int) error {
	listeningHistoryDAO := factory.ListeningHistoryDAO()

	existingHistory, err := listeningHistoryDAO.GetListeningHistory(history.ArtistId, history.UserId)
	record := &schema.ListeningHistory{
		ArtistId:    history.ArtistId,
		UserId:      history.UserId,
		Timestamp:   history.Timestamp,
		IsFollowing: history.IsFollowing,
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if history.IsFollowing {
				record.ListenCount = 1000
			} else {
				record.ListenCount = int64(counter[history.ArtistId])
			}
			if err := listeningHistoryDAO.InsertHistory(record); err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		record.ListenCount = existingHistory.ListenCount + int64(counter[history.ArtistId])
		if err := listeningHistoryDAO.UpdateHistory(history.ArtistId, history.UserId, int64(counter[history.ArtistId]), history.IsFollowing, history.Timestamp); err != nil {
			return err
		}
	}

	return nil

}

func (r *ListeningHistoryRepository) UpdateIsFollowing(factory dao.Factory, artist entity.Artist, userId string, timestamp time.Time) error {
	listeningHistoryDAO := factory.ListeningHistoryDAO()

	artistId := artist.Id
	_, err := listeningHistoryDAO.GetListeningHistory(artistId, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if err = listeningHistoryDAO.InsertHistory(&schema.ListeningHistory{
				ArtistId:    artistId,
				UserId:      userId,
				Timestamp:   timestamp,
				ListenCount: 1000,
				IsFollowing: true,
			}); err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		if err := listeningHistoryDAO.UpdateHistory(artistId, userId, 0, true, timestamp); err != nil {
			return err
		}
	}
	return nil
}
