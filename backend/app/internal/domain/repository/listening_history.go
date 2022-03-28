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
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			record := &schema.ListeningHistory{
				ArtistID:    history.ArtistID,
				UserID:      history.UserID,
				Timestamp:   history.Timestamp,
				IsFollowing: history.IsFollowing,
			}

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
		record := &schema.ListeningHistory{
			ArtistID:    existingHistory.ArtistID,
			UserID:      existingHistory.UserID,
			Timestamp:   history.Timestamp,
			IsFollowing: existingHistory.IsFollowing,
		}

		record.ListenCount = existingHistory.ListenCount + int64(count)
		if err := listeningHistoryDAO.UpdateHistory(history.ArtistID, history.UserID, int64(count), record.IsFollowing, history.Timestamp); err != nil {
			return err
		}
	}

	return nil
}

func (r *ListeningHistoryRepository) UpdateIsFollowing(factory dao.Factory, artistID, userId string, timestamp time.Time) error {
	listeningHistoryDAO := factory.ListeningHistoryDAO()

	existingHistory, err := listeningHistoryDAO.GetListeningHistory(artistID, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if err = listeningHistoryDAO.InsertHistory(&schema.ListeningHistory{
				ArtistID:    artistID,
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
		if existingHistory.IsFollowing {
			if err := listeningHistoryDAO.UpdateHistory(artistID, userId, 0, true, timestamp); err != nil {
				return err
			}
		} else {
			if err := listeningHistoryDAO.UpdateHistory(artistID, userId, 1000, true, timestamp); err != nil {
				return err
			}
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

func (r *ListeningHistoryRepository) GetListeningHistory(factory dao.Factory, artistID, userID string) (*entity.ListeningHistory, error) {
	listeningHistoryDAO := factory.ListeningHistoryDAO()
	history, err := listeningHistoryDAO.GetListeningHistory(artistID, userID)
	if err != nil {
		return nil, err
	}
	return entity.NewListeningHistory(history), nil
}
