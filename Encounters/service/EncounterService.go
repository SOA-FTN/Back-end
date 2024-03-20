package service

import (
	"encounters/repo"
)

type EncounterService struct {
	EncounterRepository *repo.EncounterRepository
}

func NewEncounterService(er *repo.EncounterRepository) *EncounterService {
	return &EncounterService{
		EncounterRepository: er,
	}
}
