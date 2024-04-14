package service

import (
	"encounters/model"
	"encounters/repo"
	"errors"
	"log"
	"time"
)

type EncounterExecutionService struct {
	logger *log.Logger
	EncounterExecutionRepository *repo.EncounterExecutionRepository
}

func NewEncounterExecutionService(l *log.Logger,er *repo.EncounterExecutionRepository) *EncounterExecutionService {
	return &EncounterExecutionService{
		l,er,
	}
}

func (es *EncounterExecutionService) CreateEncounterExecution(encounterExecution *model.EncounterExecution) error {

	newEncounterExecution := model.EncounterExecution{
		UserID:         encounterExecution.UserID,
		EncounterID:    encounterExecution.EncounterID,
		CompletionTime: encounterExecution.CompletionTime,
		IsCompleted:    encounterExecution.IsCompleted,
	}
	return es.EncounterExecutionRepository.Insert(&newEncounterExecution)
}

func (es *EncounterExecutionService) GetAllEncounterExecutions() (model.EncounterExecutions, error) {
	encounterExecutions, err := es.EncounterExecutionRepository.GetAllEncounterExecutions()
	if err != nil {
		return nil, err
	}
	return encounterExecutions, nil
}

func (ees *EncounterExecutionService) GetEncounterExecutionByUserIDAndNotCompleted(userID int) (*model.EncounterExecution, error) {
	return ees.EncounterExecutionRepository.GetByUserIDAndNotCompleted(userID);
}

func (ees *EncounterExecutionService) UpdateEncounterExecution(userID int) error {

	execution, err := ees.EncounterExecutionRepository.GetByUserIDAndNotCompleted(userID)
	if err != nil {
		return err
	}
	if execution == nil {
		return errors.New("No encounter execution found for the user")
	}

	currentTime := time.Now().String()
	execution.CompletionTime = currentTime
	execution.IsCompleted = true

	if err := ees.EncounterExecutionRepository.Update(execution.ID.Hex(),execution); err != nil {
		return err
	}
	return nil
	
}
