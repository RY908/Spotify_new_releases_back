package service

import (
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

func (s *ListeningHistoryService) InsertHistories(factory dao.Factory, artists []entity.Artist, userID string, counter map[string]int, isFollowing bool) error {
	for _, artist := range artists {
		if err := s.listeningHistoryRepository.InsertOrUpdateListeningHistory(factory,
			entity.ListeningHistory{
				ArtistID:    artist.ID,
				UserID:      userID,
				Timestamp:   time.Now(),
				IsFollowing: isFollowing,
			}, counter[artist.ID]); err != nil {
			return err
		}
	}
	return nil
}

func (s *ListeningHistoryService) UpdateIsFollowings(factory dao.Factory, artists []entity.Artist, userId string, timestamp time.Time) error {
	for _, artist := range artists {
		if err := s.listeningHistoryRepository.UpdateIsFollowing(factory, artist, userId, timestamp); err != nil {
			return err
		}
	}
	return nil
}

func (s *ListeningHistoryService) DeleteHistoriesByArtistIDs(factory dao.Factory, userID string, artistIDs []string) error {
	for _, artistID := range artistIDs {
		if err := s.listeningHistoryRepository.DeleteListeningHistoryByArtistID(factory, userID, artistID); err != nil {
			return err
		}
	}
	return nil
}

func (s *ListeningHistoryService) DeleteHistoriesByTimestamp(factory dao.Factory, userID string, timestamp time.Time) error {
	if err := s.listeningHistoryRepository.DeleteListeningHistoryByTimestamp(factory, userID, timestamp); err != nil {
		return err
	}
	return nil
}

func (s *ListeningHistoryService) GetArtistsByUserID(factory dao.Factory, userID string) (*[]entity.Artist, error) {
	artists, err := s.userArtistsRepository.GetArtistsByUserID(factory, userID)
	if err != nil {
		return nil, err
	}

	return artists, nil
}
