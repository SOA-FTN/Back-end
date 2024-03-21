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

func (er *EncounterRepository) GetAllEncounters() ([]model.Encounter, error) {
	var encounters []model.Encounter
	if err := er.DatabaseConnection.Find(&encounters).Error; err != nil {
		return nil, err
	}
	return encounters, nil
}

func (repo EncounterRepository) UpdateEncounter(encounter *model.Encounter) (*model.Encounter, error) {
	dbResult := repo.DatabaseConnection.Model(&model.Encounter{}).Where("name=?", encounter.Name).Updates(encounter)
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	return encounter, nil
}

func (er *EncounterRepository) GetEncounterByID(id int) (*model.Encounter, error) {
	var encounter model.Encounter
	if err := er.DatabaseConnection.First(&encounter, id).Error; err != nil {
		return nil, err
	}
	return &encounter, nil
}
