package service

import (
	"encounters/model"
	"encounters/repo"
)

type EncounterService struct {
	EncounterRepository *repo.EncounterRepository
}

func NewEncounterService(er *repo.EncounterRepository) *EncounterService {
	return &EncounterService{
		EncounterRepository: er,
	}
}

func (es *EncounterService) CreateEncounter(encounter *model.Encounter) error {

	newEncounter := model.Encounter{
		Name:             encounter.Name,
		Description:      encounter.Description,
		XpPoints:         encounter.XpPoints,
		Status:           encounter.Status,
		Type:             encounter.Type,
		Latitude:         encounter.Latitude,
		Longitude:        encounter.Longitude,
		ShouldBeApproved: encounter.ShouldBeApproved,
	}
	err := es.EncounterRepository.CreateEncounter(&newEncounter)
	if err != nil {
		return err
	}
	return nil
}
