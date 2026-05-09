package repository

import (
	"errors"
	"time"

	"pdf-parser/internal/model"
	"pdf-parser/pkg/database"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) Create(user *model.User) error {
	return database.GetDB().Create(user).Error
}

func (r *UserRepository) GetByUserID(userID string) (*model.User, error) {
	var user model.User
	err := database.GetDB().Where("user_id = ?", userID).First(&user).Error
	if err != nil {
		if errors.Is(err, errors.New("record not found")) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByID(id uint64) (*model.User, error) {
	var user model.User
	err := database.GetDB().First(&user, id).Error
	if err != nil {
		if errors.Is(err, errors.New("record not found")) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateUser(user *model.User) error {
	return database.GetDB().Save(user).Error
}

func (r *UserRepository) CreateOrUpdateUser(userInfo *model.User) (*model.User, error) {
	var existingUser model.User
	err := database.GetDB().Where("user_id = ?", userInfo.UserID).First(&existingUser).Error

	if err != nil && !errors.Is(err, errors.New("record not found")) {
		return nil, err
	}

	if errors.Is(err, errors.New("record not found")) {
		if err := database.GetDB().Create(userInfo).Error; err != nil {
			return nil, err
		}
		return userInfo, nil
	}

	existingUser.Username = userInfo.Username
	existingUser.Email = userInfo.Email
	existingUser.Avatar = userInfo.Avatar
	existingUser.Nickname = userInfo.Nickname
	existingUser.LastLogin = time.Now()

	if err := database.GetDB().Save(&existingUser).Error; err != nil {
		return nil, err
	}

	return &existingUser, nil
}
