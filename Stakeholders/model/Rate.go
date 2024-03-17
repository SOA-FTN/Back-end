package model

import (
	"gorm.io/gorm"
)

type Rate struct {
	gorm.Model
	Id uint `json:"Id"`
	UserId uint `json:"UserId"`
	Rating uint `json:"Rating"`
	Description string `json:"Description"`
}
