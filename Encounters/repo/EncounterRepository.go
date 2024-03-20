package repo

import (
	"encounters/model"

	"gorm.io/gorm"
)

type EncounterRepository struct {
	DatabaseConnection *gorm.DB
}

func NewEncounterRepository(db *gorm.DB) *EncounterRepository {
	return &EncounterRepository{
		DatabaseConnection: db,
	}
}

func (tr *EncounterRepository) CreateEncounter(encounter *model.Encounter) error {
	dbResult := tr.DatabaseConnection.Create(encounter)
	if dbResult.Error != nil {
		return dbResult.Error
	}
	return nil
}
