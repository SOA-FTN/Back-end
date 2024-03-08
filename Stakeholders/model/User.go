package model

import (
	"gorm.io/gorm"
)


type UserRole int

const(
	Admin UserRole = iota
	Tourist
	Author
)

type User struct {
	gorm.Model
	UserName string `json:"username"`
	Password string `json:"password"`
	Role UserRole `json:"role"`
	IsActive bool `json:"isActive"`
	Person Person `gorm:"foreignKey:UserID"`
}
