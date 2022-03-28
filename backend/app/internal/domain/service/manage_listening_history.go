package service

import (
	"fmt"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/entity"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/repository"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/dao"
	"time"
)

func NewListeningHistoryService() *ListeningHistoryService {
	listeningHistoryRepository := repository.NewListeningHistoryRepository()
	userArtistsRepository := repository.NewUserArtistsRepository()

	return &ListeningHistoryService{
		listeningHistoryRepository: listeningHistoryRepository,
		userArtistsRepository:      userArtistsRepository,
	}
}

type ListeningHistoryService struct {
	listeningHistoryRepository *repository.ListeningHistoryRepository
	userArtistsRepository      *repository.UserArtistsRepository
}

func (s *ListeningHistoryService) InsertHistories(factory dao.Factory, artists []*entity.Artist, userID string, counter map[string]int, isFollowing bool, timestamp time.Time) error {
	for _, artist := range artists {
		if err := s.listeningHistoryRepository.InsertOrUpdateListeningHistory(factory,
			entity.ListeningHistory{
				ArtistID:    artist.ID,
				UserID:      userID,
				Timestamp:   timestamp,
				IsFollowing: isFollowing,
			}, counter[artist.ID]); err != nil {
			return fmt.Errorf("unable to insert listening histories: %w", err)
		}
	}
	return nil
}

func (s *ListeningHistoryService) UpdateIsFollowings(factory dao.Factory, artists []*entity.Artist, userId string, timestamp time.Time) error {
	for _, artist := range artists {
		if err := s.listeningHistoryRepository.UpdateIsFollowing(factory, artist.ID, userId, timestamp); err != nil {
			return fmt.Errorf("unable to update isFollowing in listening history: %w", err)
		}
	}
	return nil
}

func (s *ListeningHistoryService) DeleteHistoriesByArtistIDs(factory dao.Factory, userID string, artistIDs []string) error {
	for _, artistID := range artistIDs {
		if err := s.listeningHistoryRepository.DeleteListeningHistoryByArtistID(factory, userID, artistID); err != nil {
			return fmt.Errorf("unable to delete listening histories by artist id: %w", err)
		}
	}
	return nil
}

func (s *ListeningHistoryService) DeleteHistoriesByTimestamp(factory dao.Factory, userID string, timestamp time.Time) error {
	if err := s.listeningHistoryRepository.DeleteListeningHistoryByTimestamp(factory, userID, timestamp); err != nil {
		return fmt.Errorf("unable to delete listening histories by timestamp: %w", err)
	}
	return nil
}

func (s *ListeningHistoryService) GetArtistsByUserID(factory dao.Factory, userID string) ([]*entity.UserArtist, error) {
	artists, err := s.userArtistsRepository.GetArtistsByUserID(factory, userID)
	if err != nil {
		return nil, fmt.Errorf("unable get artists by user id: %w", err)
	}

	return artists, nil
}
