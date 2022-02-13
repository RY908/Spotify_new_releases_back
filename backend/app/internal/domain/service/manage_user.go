package service

import (
	"fmt"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/entity"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/repository"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/dao"
	"log"
	"time"
)

func NewUserService() *UserService {
	userRepository := repository.NewUserRepository()
	return &UserService{
		repository: userRepository,
	}
}

type UserService struct {
	logger     *log.Logger
	repository *repository.UserRepository
}

type InsertUserParam struct {
	UserId        string
	AccessToken   string
	TokenType     string
	RefreshToken  string
	Expiry        time.Time
	PlaylistId    string
	IfRemixAdd    bool
	IfAcousticAdd bool
}

func (s *UserService) InsertUser(factory dao.Factory, user entity.User) error {
	if err := s.repository.InsertUser(factory, user); err != nil {
		return fmt.Errorf("unable to insert user: %w", err)
	}
	return nil
}

func (s *UserService) GetUser(factory dao.Factory, userID string) (*entity.User, error) {
	user, err := s.repository.GetUserByUID(factory, userID)
	if err != nil {
		return nil, fmt.Errorf("unable to get user: %w", err)
	}
	return user, nil
}

func (s *UserService) GetAllUsers(factory dao.Factory) ([]*entity.User, error) {
	users, err := s.repository.GetAllUsers(factory)
	if err != nil {
		return nil, fmt.Errorf("unable to get all users: %w", err)
	}
	return users, nil
}

func (s *UserService) UpdateUserToken(factory dao.Factory, user entity.User) error {
	if err := s.repository.UpdateUserToken(factory, user); err != nil {
		return fmt.Errorf("unable to update user token: %w", err)
	}
	return nil
}

func (s *UserService) UpdateUserPreference(factory dao.Factory, user entity.User) error {
	if err := s.repository.UpdateUserPreference(factory, user); err != nil {
		return fmt.Errorf("unable to update user preference: %w", err)
	}
	return nil
}
