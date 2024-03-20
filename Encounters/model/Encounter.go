package model

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
	Name             string
	Description      string
	XpPoints         int
	Status           EncounterStatus
	Type             EncounterType
	Latitude         float64
	Longitude        float64
	ShouldBeApproved bool
}
