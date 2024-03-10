package service

import (
	"fmt"
	"stakeholders/model"
	"stakeholders/repo"
)

type UserService struct {
	UserRepo *repo.UserRepository
}
//GETOVANJE PROFILA
func (service *UserService) GetPersonByUserId(userId *string) (*model.Person , error){
	person,err := service.UserRepo.GetPersonByUserId(userId)
	if err != nil {
		return nil , fmt.Errorf(fmt.Sprintf("menu item with id %s not found", *userId))
	}
	return person,nil
}
//REGISTRACIJA
func(service *UserService) RegisterUser (user *model.User, person *model.Person) error {

	err := service.UserRepo.RegisterUser(user)
	if err != nil{
		return err
	}

	person.UserID = user.ID

	err = service.UserRepo.RegisterPerson(person)
	if err != nil {
		return err
	}
	return nil
}
//UPDATE
func (service *UserService) UpdateProfile(person *model.Person) (*model.Person,error){
	updatedPerson , err := service.UserRepo.UpdateProfile(person)
	if err != nil {
		return nil,err;
	}
	return updatedPerson,nil
}