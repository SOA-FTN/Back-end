package repo

import (
	"stakeholders/model"

	"gorm.io/gorm"
)

type AuthRepository struct {
	DatabaseConnection *gorm.DB
}

func (authRepo *AuthRepository) Authentication(credentials *model.Credentials) (*model.User, error) {
	var user model.User
	dbResult := authRepo.DatabaseConnection.Where("user_name = ? AND password = ?", credentials.Username, credentials.Password).First(&user)
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	return &user, nil
}
