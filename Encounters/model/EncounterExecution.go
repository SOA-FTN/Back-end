package model

import (
	"time"

	"gorm.io/gorm"
)

type EncounterExecution struct {
	gorm.Model
	UserID         int64
	EncounterID    int64
	CompletionTime *time.Time
	IsCompleted    bool
}
