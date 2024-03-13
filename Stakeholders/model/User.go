package model

import (
	"gorm.io/gorm"
)


type UserRole int

const(
	admin UserRole = iota
	tourist
	author
)

type User struct {
	gorm.Model
	UserName string `json:"username"`
	Password string `json:"password"`
	Role UserRole `json:"role"`
	IsActive bool `json:"isActive"`
	Person Person `gorm:"foreignKey:UserID"`
}


func (u *User) GetRoleName() string {
    switch u.Role {
    case admin:
        return "admin"
    case tourist:
        return "tourist"
    case author:
        return "author"
    default:
        return ""
    }
}
