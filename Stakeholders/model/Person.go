package model

import (
	"gorm.io/gorm"
)

type Person struct {
	gorm.Model
	UserID uint `json:"userId"`
	Name string `json:"name"`
	Surname string `json:"surname"`
	Email string `json:"email"`
	Bio string `json:"bio"`
	Quote string `json:"quote"`
}


