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

func (r *ListeningHistoryRepository) InsertOrUpdateListeningHistory(factory dao.Factory, history entity.ListeningHistory, count int) error {
	listeningHistoryDAO := factory.ListeningHistoryDAO()

	existingHistory, err := listeningHistoryDAO.GetListeningHistory(history.ArtistID, history.UserID)
	record := &schema.ListeningHistory{
		ArtistID:    history.ArtistID,
		UserID:      history.UserID,
		Timestamp:   history.Timestamp,
		IsFollowing: history.IsFollowing,
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if history.IsFollowing {
				record.ListenCount = 1000
			} else {
				record.ListenCount = int64(count)
			}
			if err := listeningHistoryDAO.InsertHistory(record); err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		record.ListenCount = existingHistory.ListenCount + int64(count)
		if err := listeningHistoryDAO.UpdateHistory(history.ArtistID, history.UserID, int64(count), history.IsFollowing, history.Timestamp); err != nil {
			return err
		}
	}

	return nil
}

func (r *ListeningHistoryRepository) UpdateIsFollowing(factory dao.Factory, artist entity.Artist, userId string, timestamp time.Time) error {
	listeningHistoryDAO := factory.ListeningHistoryDAO()

	artistId := artist.ID
	_, err := listeningHistoryDAO.GetListeningHistory(artistId, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if err = listeningHistoryDAO.InsertHistory(&schema.ListeningHistory{
				ArtistID:    artistId,
				UserID:      userId,
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

func (r *ListeningHistoryRepository) DeleteListeningHistoryByArtistID(factory dao.Factory, userID, artistID string) error {
	listeningHistoryDAO := factory.ListeningHistoryDAO()
	if err := listeningHistoryDAO.DeleteHistoryByArtistId(userID, artistID); err != nil {
		return err
	}
	return nil
}

func (r *ListeningHistoryRepository) DeleteListeningHistoryByTimestamp(factory dao.Factory, userID string, timestamp time.Time) error {
	listeningHistoryDAO := factory.ListeningHistoryDAO()
	if err := listeningHistoryDAO.DeleteHistoryByTimestamp(userID, timestamp); err != nil {
		return err
	}
	return nil
}
