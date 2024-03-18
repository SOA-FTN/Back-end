package model

import (
	"gorm.io/gorm"
)

type Person struct {
	gorm.Model
	UserID uint `json:"UserId"`
	Name string `json:"Name"`
	Surname string `json:"Surname"`
	Image string `json:"ProfileImage"`
	Email string `json:"Email"`
	Bio string `json:"Bio"`
	Quote string `json:"Quote"`
}


