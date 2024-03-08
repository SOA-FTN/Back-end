package repo

import (
	"stakeholders/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	DatabaseConnection *gorm.DB
}

func(userRepo *UserRepository) RegisterUser(user *model.User) error {
	dbResult := userRepo.DatabaseConnection.Create(user)
	if(dbResult.Error != nil) {
		return dbResult.Error
	}
	return nil
}

func(userRepo *UserRepository) RegisterPerson(person *model.Person) error {
	dbResult := userRepo.DatabaseConnection.Create(person)
	if dbResult.Error != nil {
		return dbResult.Error
	}
	return nil
}

func (userRepo *UserRepository) Authentication(credentials *model.Credentials) (*model.User,error) {
	var user model.User
	dbResult := userRepo.DatabaseConnection.Where("user_name = ? AND password = ?",credentials.Username,credentials.Password).First(&user)
	if dbResult.Error != nil {
		return nil,dbResult.Error
	}
	return &user,nil
}
