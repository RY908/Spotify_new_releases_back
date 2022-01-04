package repository

import (
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/entity"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/dao"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/schema"
)

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

type UserRepository struct{}

func (r *UserRepository) InsertUser(factory dao.Factory, user entity.User) error {
	userDAO := factory.UserDAO()
	record := &schema.User{
		Id:            user.Id,
		AccessToken:   user.AccessToken,
		TokenType:     user.TokenType,
		RefreshToken:  user.RefreshToken,
		Expiry:        user.Expiry,
		PlaylistId:    user.PlaylistId,
		IfRemixAdd:    user.IfRemixAdd,
		IfAcousticAdd: user.IfAcousticAdd,
	}
	if err := userDAO.InsertUser(record); err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) UpdateUserToken(factory dao.Factory, user entity.User) error {
	userDAO := factory.UserDAO()

	existingUser, err := userDAO.GetUser(user.Id)
	if err != nil {
		return err
	}

	existingUser.AccessToken = user.AccessToken
	existingUser.TokenType = user.TokenType
	existingUser.RefreshToken = user.RefreshToken
	existingUser.Expiry = user.Expiry
	if err := userDAO.UpdateUser(existingUser); err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) UpdateUserPreference(factory dao.Factory, user entity.User) error {
	userDAO := factory.UserDAO()

	existingUser, err := userDAO.GetUser(user.Id)
	if err != nil {
		return err
	}
	existingUser.IfRemixAdd = user.IfRemixAdd
	existingUser.IfAcousticAdd = user.IfAcousticAdd
	if err := userDAO.UpdateUser(existingUser); err != nil {
		return err
	}
	return nil
}