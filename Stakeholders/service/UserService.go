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
func(service *UserService) Registration (registration *model.Registration, token *string) error {

	newUser := model.User{
		UserName: registration.Username,
		Password: registration.Password,
		IsActive: false,
		Role:model.ParseUserRole(registration.Role),
		VerificationToken: *token,
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

//GETOVANJE PROFILA
func (service *UserService) GetAndVerifyUserByToken(token *string) (*model.User , error){
	user,err := service.UserRepo.GetUserByToken(token)
	if err != nil {
		return nil , fmt.Errorf(fmt.Sprintf("menu item with token %s not found", *token))
	}

	user.IsActive = true
	updatedUser, err := service.UserRepo.UpdateUser(user)
	if err != nil {
		return nil,err;
	}

	return updatedUser,nil
}

