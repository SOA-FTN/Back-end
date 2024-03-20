package repo

import (
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
