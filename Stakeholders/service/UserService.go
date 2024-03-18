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
func(service *UserService) Registration (registration *model.Registration) error {

	newUser := model.User{
		UserName: registration.Username,
		Password: registration.Password,
		IsActive: true,
		Role:model.ParseUserRole(registration.Role),
	}
	newPerson := model.Person{
		Name:    registration.Name,
		Surname: registration.Surname,
		Email:   registration.Email,
	}

	err := service.UserRepo.RegisterUser(&newUser)
	if err != nil{
		return err
	}

	newPerson.UserID = newUser.ID

	err = service.UserRepo.RegisterPerson(&newPerson)
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