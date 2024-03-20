package handler

import (
	"encoding/json"
	"encounters/model"
	"encounters/service"
	"log"
	"net/http"
)

type EncounterHandler struct {
	EncounterService *service.EncounterService
}

func NewEncounterHandler(es *service.EncounterService) *EncounterHandler {
	return &EncounterHandler{
		EncounterService: es,
	}
}

func (eh *EncounterHandler) CreateEncounterHandler(w http.ResponseWriter, r *http.Request) {

	var enc model.CreateEncounter
	if err := json.NewDecoder(r.Body).Decode(&enc); err != nil {
		log.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Map fields to the Tour struct

	newEncounter := model.Encounter{
		Name:             enc.Name,
		Description:      enc.Description,
		XpPoints:         enc.XpPoints,
		Status:           model.EncounterStatus(service.ConvertEncounterStatusToInt(enc.Status)),
		Type:             model.EncounterType(service.ConvertEncounterTypeToInt(enc.Type)),
		Latitude:         enc.Latitude,
		Longitude:        enc.Longitude,
		ShouldBeApproved: enc.ShouldBeApproved,
	}
	if err := eh.EncounterService.CreateEncounter(&newEncounter); err != nil {
		http.Error(w, "Failed to create encounter", http.StatusInternalServerError)
		log.Println("ne")
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(enc)
}
