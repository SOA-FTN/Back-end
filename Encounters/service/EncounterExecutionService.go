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

func (es *EncounterExecutionService) GetAllEncounterExecutions() ([]model.EncounterExecution, error) {
	encounters, err := es.EncounterExecutionRepository.GetAllEncounterExecutions()
	if err != nil {
		return nil, err
	}
	return encounters, nil
}
