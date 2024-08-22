package model

import "gorm.io/gorm"

type TourPointType int

const (
	Start TourPointType = iota
	End
	InBetween
)

type TourPoint struct {
	gorm.Model
	TourId      int64         `json:"tourId"`
	Name        string        `json:"name"`
	Description *string       `json:"description,omitempty"`
	Latitude    float64       `json:"latitude"`
	Longitude   float64       `json:"longitude"`
	ImageUrl    string        `json:"imageUrl"`
	Type        TourPointType `json:"type"`
}
