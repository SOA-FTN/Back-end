package handler

import (
	"encounters/service"
)

type EncounterHandler struct {
	EncounterService *service.EncounterService
}

func NewEncounterHandler(es *service.EncounterService) *EncounterHandler {
	return &EncounterHandler{
		EncounterService: es,
	}
}
