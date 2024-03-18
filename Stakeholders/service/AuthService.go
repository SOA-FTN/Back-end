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

var jwtKey = string ("explorer_secret_key")

//Authentifikacija
func (service *AuthService) Authentication (credentials *model.Credentials) (*model.User,error) {
	user,err := service.AuthRepo.Authentication(credentials)
	return user,err
}
//Generisanje tokena
func (service *AuthService) GenerateToken(user *model.User) (string, error) {
    expirationTime := time.Now().Add(time.Minute * 60 * 24) 

    claims := jwt.MapClaims{
        "id":       user.ID,
        "username": user.UserName,
        "role":     user.GetRoleName(),
        "exp":      expirationTime.Unix(),
        "iss":      "explorer",
        "aud":      "explorer-front.com",
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(jwtKey))
    if err != nil {
        return "", err
    }

    return tokenString, nil
}