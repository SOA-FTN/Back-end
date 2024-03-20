package model

import "time"

type EncounterExecution struct {
	UserID         int64      // Assuming the user ID is of type long in C#
	EncounterID    int64      // Assuming the encounter ID is of type long in C#
	CompletionTime *time.Time // Nullable DateTime equivalent in Go
	IsCompleted    bool
}
