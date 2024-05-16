package model

import (
	"encoding/json"
	"io"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
	ID           	 primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	Name             string				 `bson:"name" json:"name"`
	Description      string 			 `bson:"description,omitempty" json:"description"`
	XpPoints         int32				 `bson:"xppoints" json:"xppoints"`
	Status           EncounterStatus	 `bson:"status" json:"status"`
	Type             EncounterType		 `bson:"type" json:"type"`
	Latitude         float64  			 `bson:"latitude,omitempty" json:"latitude"`
	Longitude        float64			 `bson:"longitude,omitempty" json:"longitude"`
	ShouldBeApproved bool				 `bson:"shouldbeapproved" json:"shouldbeapproved"`
}

type CreateEncounter struct {
	Name             string				`bson:"name,omitempty" json:"name"`
	Description      string				`bson:"description,omitempty" json:"description"`
	XpPoints         int32				`bson:"xppoints" json:"xppoints"`
	Status           string				`bson:"status" json:"status"`
	Type             string				`bson:"type" json:"type"`
	Latitude         float64			`bson:"latitude,omitempty" json:"latitude"`
	Longitude        float64			`bson:"longitude,omitempty" json:"longitude"`
	ShouldBeApproved bool				`bson:"shouldbeapproved" json:"shouldbeapproved"`
}

type UpdateEncounter struct {
	Name             string				`bson:"name,omitempty" json:"name"`
	Description      string				`bson:"description,omitempty" json:"description"`
	XpPoints         int				`bson:"points,omitempty" json:"points"`
	Status           EncounterStatus	`bson:"status,omitempty" json:"status"`
	Type             string				`bson:"type,omitempty" json:"type"`
	Latitude         float64			`bson:"latitude,omitempty" json:"latitude"`
	Longitude        float64			`bson:"longitude,omitempty" json:"longitude"`
	ShouldBeApproved bool				`bson:"approved,omitempty" json:"approved"`
}

type Encounters []*Encounter

func (p *Encounters) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (encounter *Encounter) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(encounter);
}

func (encounter *Encounter) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(encounter)
}


func(encounter *CreateEncounter) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(encounter);
}

func (encounter *CreateEncounter) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(encounter)
}

