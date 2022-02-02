package service

import (
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/entity"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/repository"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/dao"
	"time"
)

func NewUserService() *UserService {
	userRepository := repository.NewUserRepository()
	return &UserService{
		repository: userRepository,
	}
}

type UserService struct {
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
		return err
	}
	return nil
}

func (s *UserService) GetUser(factory dao.Factory, userID string) (*entity.User, error) {
	user, err := s.repository.GetUserByUID(factory, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetAllUsers(factory dao.Factory) ([]*entity.User, error) {
	users, err := s.repository.GetAllUsers(factory)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) UpdateUserToken(factory dao.Factory, user entity.User) error {
	if err := s.repository.UpdateUserToken(factory, user); err != nil {
		return err
	}
	return nil
}

func (s *UserService) UpdateUserPreference(factory dao.Factory, user entity.User) error {
	if err := s.repository.UpdateUserPreference(factory, user); err != nil {
		return err
	}
	return nil
}
