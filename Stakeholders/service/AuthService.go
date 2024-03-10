package service

import (
	"stakeholders/model"
	"stakeholders/repo"
	"time"

	"github.com/golang-jwt/jwt"
)

type AuthService struct {
	AuthRepo *repo.AuthRepository
}

var jwtKey = []byte ("secret_key")

//Authentifikacija
func (service *AuthService) Authentication (credentials *model.Credentials) (*model.User,error) {
	user,err := service.AuthRepo.Authentication(credentials)
	return user,err
}
//Generisanje tokena
func (authService *AuthService) GenerateToken(user *model.User) (string,time.Time, error) {
	expirationTime := time.Now().Add(time.Minute *5)
	claims := authService.CreateClaims(user,&expirationTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, *claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", time.Time{}, err
	}
	return tokenString, expirationTime, nil
}

//Kreiranje Claim-a tokena
func (service *AuthService) CreateClaims(user *model.User,expirationTime *time.Time) *model.Claims {
	claims := model.Claims{
		Username: user.UserName,
		Role: user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	return &claims
}