package service

import (
	"encounters/model"
	"encounters/repo"
	"errors"
	"log"
)

type EncounterService struct {
	logger *log.Logger
	EncounterRepository *repo.EncounterRepository
}

func NewEncounterService(l *log.Logger,er *repo.EncounterRepository) *EncounterService {
	return &EncounterService{
		l,er,
	}
}

func (es *EncounterService) CreateEncounter(encounter *model.CreateEncounter) error {
	if es.EncounterRepository == nil {
		log.Println("REPOSITORY ERRORCINA")
        return errors.New("EncounterRepository is nil")
    }
	newEncounter := model.Encounter{
		Name:             encounter.Name,
		Description:      encounter.Description,
		XpPoints:         encounter.XpPoints,
		Status:           model.EncounterStatus(ConvertEncounterStatusToInt(encounter.Status)),
		Type:             model.EncounterType(ConvertEncounterTypeToInt(encounter.Type)),
		Latitude:         encounter.Latitude,
		Longitude:        encounter.Longitude,
		ShouldBeApproved: encounter.ShouldBeApproved,
	}
	log.Println("=========================")
	log.Println(newEncounter)
	return es.EncounterRepository.Insert(&newEncounter)
}


/*
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
*/

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


func (es *EncounterService) GetAllEncounters() (model.Encounters, error) {
	encounters, err := es.EncounterRepository.GetAllEncounters()
	if err != nil {
		return nil, err
	}
	return encounters, nil
}
/*
func (service *EncounterService) UpdateEncounter(encounter *model.Encounter) (*model.Encounter, error) {
	updatedEncounter, err := service.EncounterRepository.UpdateEncounter(encounter)
	if err != nil {
		return nil, err
	}
	return updatedEncounter, nil
}
*/
func (es *EncounterService) GetEncounterByID(id string) (*model.Encounter, error) {
	return es.EncounterRepository.GetEncounterByID(id)
}

