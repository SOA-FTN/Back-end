package repo

import (
	"encounters/model"

	"gorm.io/gorm"
)

type EncounterExecutionRepository struct {
	DatabaseConnection *gorm.DB
}

func NewEncounterExecutionRepository(db *gorm.DB) *EncounterExecutionRepository {
	return &EncounterExecutionRepository{
		DatabaseConnection: db,
	}
}

func (tr *EncounterExecutionRepository) CreateEncounterExecution(encounterExecution *model.EncounterExecution) error {
	dbResult := tr.DatabaseConnection.Create(encounterExecution)
	if dbResult.Error != nil {
		return dbResult.Error
	}
	return nil
}

func (er *EncounterExecutionRepository) GetAllEncounterExecutions() ([]model.EncounterExecution, error) {
	var encounterExecutions []model.EncounterExecution
	if err := er.DatabaseConnection.Find(&encounterExecutions).Error; err != nil {
		return nil, err
	}
	return encounterExecutions, nil
}

func (eer *EncounterExecutionRepository) GetByUserIDAndNotCompleted(userID int) (*model.EncounterExecution, error) {
	var execution model.EncounterExecution
	if err := eer.DatabaseConnection.Where("user_id = ? AND is_completed = ?", userID, false).First(&execution).Error; err != nil {
		return nil, err
	}
	return &execution, nil
}
