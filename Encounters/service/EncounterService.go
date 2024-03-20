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

func ConvertEncounterStatusToInt(status string) int {
	switch status {
	case "ACTIVE":
		return 0
	case "DRAFT":
		return 1
	case "ARCHIVED":
		return 2
	default:
		return -1
	}
}

func ConvertEncounterTypeToInt(encounterType string) int {
	switch encounterType {
	case "SOCIAL":
		return 0
	case "LOCATION":
		return 1
	case "MISC":
		return 2
	default:
		return -1
	}
}

func (es *EncounterService) GetAllEncounters() ([]model.Encounter, error) {
	encounters, err := es.EncounterRepository.GetAllEncounters()
	if err != nil {
		return nil, err
	}
	return encounters, nil
}

func (service *EncounterService) UpdateEncounter(encounter *model.Encounter) (*model.Encounter, error) {
	updatedEncounter, err := service.EncounterRepository.UpdateEncounter(encounter)
	if err != nil {
		return nil, err
	}
	return updatedEncounter, nil
}
