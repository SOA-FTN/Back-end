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

func(repo UserRepository) UpdateProfile(person *model.Person) (*model.Person,error) {
	dbResult :=repo.DatabaseConnection.Model(&model.Person{}).Where("user_id=?",person.UserID).Updates(person)
	if dbResult.Error != nil {
		return nil,dbResult.Error
	}
	return person,nil
}

func(repo UserRepository) UpdateUser(user *model.User) (*model.User,error) {
	dbResult :=repo.DatabaseConnection.Model(&model.User{}).Where("id=?",user.ID).Updates(user)
	if dbResult.Error != nil {
		return nil,dbResult.Error
	}
	return user,nil
}


func(repo *UserRepository) GetPersonByUserId(userId *string) (*model.Person,error) {
	person := model.Person{}
	dbResult := repo.DatabaseConnection.First(&person, "user_id = ?", *userId)
	if dbResult.Error !=nil {
		return nil,dbResult.Error
	}
	return &person,nil
}

func(repo *UserRepository) GetUserByToken(token *string) (*model.User,error) {
	user := model.User{}
	dbResult := repo.DatabaseConnection.First(&user, "verification_token = ?", *token)
	if dbResult.Error !=nil {
		return nil,dbResult.Error
	}
	return &user,nil
}
