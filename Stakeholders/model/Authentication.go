package model

import (
	"github.com/golang-jwt/jwt"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	Role UserRole `json:"role"`
	jwt.StandardClaims
}