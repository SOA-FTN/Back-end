package service

import (
	"fmt"
	"stakeholders/DtoObjects"
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

//GETOVANJE PROFILA
func (service *UserService) GetUserByUserId(userId *uint) (*model.User , error){
	user,err := service.UserRepo.GetUserById(userId)
	if err != nil {
		return nil , fmt.Errorf(fmt.Sprintf("menu item with id %d not found", *userId))
	}
	return user,nil
}

//REGISTRACIJA
func(service *UserService) Registration (registration *model.Registration, token *string, item *bool ) error {

	newUser := model.User{
		UserName: registration.Username,
		Password: registration.Password,
		Role:model.ParseUserRole(registration.Role),
		VerificationToken: *token,
		IsActive: item,
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

func (service *UserService) UpdateUser(user *model.User) (*model.User,error){
	updatedUser , err := service.UserRepo.UpdateUser(user)
	if err != nil {
		return nil,err;
	}
	return updatedUser,nil
}


func (service *UserService) GetAndVerifyUserByToken(token *string) (*model.User , error){
	user,err := service.UserRepo.GetUserByToken(token)
	if err != nil {
		return nil , fmt.Errorf(fmt.Sprintf("menu item with token %s not found", *token))
	}
	*user.IsActive = true
	updatedUser, err := service.UpdateUser(user)
	if err != nil {
		return nil,err;
	}

	return updatedUser,nil
}

func (service *UserService) GetAllProfiles() ([]DtoObjects.ProfileDto, error){
	profiles,err := service.UserRepo.GetAllProfiles()
	if err != nil {
		return nil , err
	}
	return profiles,nil
}

func (service *UserService) BlockOrUnblock(id *uint) (*model.User, error) {
    // Get the user from the repository
    user, err := service.GetUserByUserId(id)
    if err != nil {
        return nil, err // Returning error if GetUserById returns an error
    }
    
    // Toggle the IsActive status of the user
	//fmt.Println(user)
    *user.IsActive = !*user.IsActive 
    var updatedUser *model.User
    // Update the user in the repository
	//fmt.Println(user)
    updatedUser, err = service.UserRepo.UpdateUser(user)
    if err != nil {
        return nil, err // Returning error if Update returns an error
    }
	//fmt.Println(updatedUser)
    return updatedUser, nil
}

