package service

import (
	"encounters/model"
	"encounters/repo"
	"errors"
	"time"
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

func (ees *EncounterExecutionService) GetEncounterExecutionByUserIDAndNotCompleted(userID int) (*model.EncounterExecution, error) {

	execution, err := ees.EncounterExecutionRepository.GetByUserIDAndNotCompleted(userID)
	if err != nil {
		return nil, err
	}
	return execution, nil
}

func (ees *EncounterExecutionService) UpdateEncounterExecution(userID int) error {

	execution, err := ees.EncounterExecutionRepository.GetByUserIDAndNotCompleted(userID)
	if err != nil {
		return err
	}
	if execution == nil {
		return errors.New("No encounter execution found for the user")
	}

	currentTime := time.Now()
	execution.CompletionTime = &currentTime
	execution.IsCompleted = true

	if err := ees.EncounterExecutionRepository.Update(execution); err != nil {
		return err
	}
	return nil
}
