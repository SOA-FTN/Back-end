package service

import (
	"encounters/model"
	"encounters/repo"
)

type EncounterExecutionService struct {
	EncounterExecutionRepository *repo.EncounterExecutionRepository
}

func NewEncounterExecutionService(er *repo.EncounterExecutionRepository) *EncounterExecutionService {
	return &EncounterExecutionService{
		EncounterExecutionRepository: er,
	}
}

func (es *EncounterExecutionService) CreateEncounterExecution(encounterExecution *model.EncounterExecution) error {

	newEncounterExecution := model.EncounterExecution{
		UserID:         encounterExecution.UserID,
		EncounterID:    encounterExecution.EncounterID,
		CompletionTime: encounterExecution.CompletionTime,
		IsCompleted:    encounterExecution.IsCompleted,
	}
	err := es.EncounterExecutionRepository.CreateEncounterExecution(&newEncounterExecution)
	if err != nil {
		return err
	}
	return nil
}

func (es *EncounterExecutionService) GetAllEncounterExecutions() ([]model.EncounterExecution, error) {
	encounters, err := es.EncounterExecutionRepository.GetAllEncounterExecutions()
	if err != nil {
		return nil, err
	}
	return encounters, nil
}
