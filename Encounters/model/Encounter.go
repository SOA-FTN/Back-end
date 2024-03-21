package model

import "gorm.io/gorm"

type EncounterStatus int

const (
	ACTIVE EncounterStatus = iota
	DRAFT
	ARCHIVED
)

type EncounterType int

const (
	SOCIAL EncounterType = iota
	LOCATION
	MISC
)

type Encounter struct {
	gorm.Model
	Name             string
	Description      string
	XpPoints         int
	Status           EncounterStatus
	Type             EncounterType
	Latitude         float64
	Longitude        float64
	ShouldBeApproved bool
}

type CreateEncounter struct {
	Name             string
	Description      string
	XpPoints         int
	Status           string
	Type             string
	Latitude         float64
	Longitude        float64
	ShouldBeApproved bool
}

type UpdateEncounter struct {
	Name             string
	Description      string
	XpPoints         int
	Status           EncounterStatus
	Type             string
	Latitude         float64
	Longitude        float64
	ShouldBeApproved bool
}
