package service

import (
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
