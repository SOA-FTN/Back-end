package service

import (
	"stakeholders/model"
	"stakeholders/repo"
	"time"

	"github.com/golang-jwt/jwt"
)

type UserService struct {
	UserRepo *repo.UserRepository
}

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

func (service *UserService) Authentication (credentials *model.Credentials) (*model.User,error) {
	user,err := service.UserRepo.Authentication(credentials)
	return user,err
}

func (service *UserService) CreateClaims(user *model.User) *model.Claims {
	expirationTime := time.Now().Add(time.Minute *5)
	claims := model.Claims{
		Username: user.UserName,
		Role: user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	return &claims
}