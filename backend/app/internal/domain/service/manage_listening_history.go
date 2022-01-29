package service

import (
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/entity"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/repository"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/dao"
	"time"
)

func NewListeningHistoryService() *ListeningHistoryService {
	listeningHistoryRepository := repository.NewListeningHistoryRepository()
	return &ListeningHistoryService{
		listeningHistoryRepository: listeningHistoryRepository,
	}
}

type ListeningHistoryService struct {
	listeningHistoryRepository *repository.ListeningHistoryRepository
}

func (r *ListeningHistoryService) InsertHistories(factory dao.Factory, histories []entity.ListeningHistory, counter map[string]int) error {
	for _, history := range histories {
		if err := r.listeningHistoryRepository.InsertOrUpdateListeningHistory(factory, history, counter[history.ArtistId]); err != nil {
			return err
		}
	}
	return nil
}

func (r *ListeningHistoryService) UpdateIsFollowings(factory dao.Factory, artists []entity.Artist, userId string, timestamp time.Time) error {
	for _, artist := range artists {
		if err := r.listeningHistoryRepository.UpdateIsFollowing(factory, artist, userId, timestamp); err != nil {
			return err
		}
	}
	return nil
}

func (r *ListeningHistoryService) DeleteHistoriesByArtistIDs(factory dao.Factory, userID string, artistIDs []string) error {
	for _, artistID := range artistIDs {
		if err := r.listeningHistoryRepository.DeleteListeningHistoryByArtistID(factory, userID, artistID); err != nil {
			return err
		}
	}
	return nil
}

func (r *ListeningHistoryService) DeleteHistoriesByTimestamp(factory dao.Factory, userID string, timestamp time.Time) error {
	if err := r.listeningHistoryRepository.DeleteListeningHistoryByTimestamp(factory, userID, timestamp); err != nil {
		return err
	}
	return nil
}
