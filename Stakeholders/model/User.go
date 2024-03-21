package model

import (
	"gorm.io/gorm"
)


type UserRole int

const(
	administrator UserRole = iota
	tourist
	author
)

type User struct {
	gorm.Model
	UserName string `json:"Username"`
	Password string `json:"Password"`
	Role UserRole `json:"Role"`
	IsActive *bool `json:"IsActive"`
	VerificationToken string `json:"VerificationToken"`
	Person Person `gorm:"foreignKey:UserID"`
}



func (u *User) GetRoleName() string {
	switch u.Role {
	case administrator:
        return "administrator"
    case tourist:
        return "tourist"
    case author:
        return "author"
    default:
        return ""
    }
}

func ParseUserRole(role string) UserRole {
	switch role {
	case "administrator":
		return administrator
	case "tourist":
		return tourist
	case "author":
		return author
	default:
		return tourist // Defaultna vrijednost, možete promijeniti prema potrebi
	}
}

type Registration struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
	Email string `json:"Email"`
	Name string `json:"Name"`
	Surname string `json:"Surname"`
	Role string `json:"Role"`
}

