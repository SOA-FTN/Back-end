package handler

import (
	"encoding/json"
	"encounters/model"
	"encounters/service"
	"strconv"

	"log"
	"net/http"

	"github.com/gorilla/mux"
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

func (eh *EncounterHandler) GetAllEncountersHandler(w http.ResponseWriter, r *http.Request) {
	encounters, err := eh.EncounterService.GetAllEncounters()
	if err != nil {
		http.Error(w, "Failed to get encounters", http.StatusInternalServerError)
		log.Println("prva")
		return
	}

	// Convert encounters to JSON and send response
	response, err := json.Marshal(encounters)
	if err != nil {
		http.Error(w, "Failed to marshal encounters", http.StatusInternalServerError)
		log.Println("druga")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (eh *EncounterHandler) GetEncounterByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Extract encounter ID from the URL path
	vars := mux.Vars(r)
	encounterIDStr, ok := vars["encounterId"]
	if !ok {
		http.Error(w, "Encounter ID not provided", http.StatusBadRequest)
		return
	}

	// Convert encounterIDStr to int
	encounterID, err := strconv.Atoi(encounterIDStr)
	if err != nil {
		http.Error(w, "Invalid encounter ID", http.StatusBadRequest)
		return
	}

	// Call service function to get encounter by ID
	encounter, err := eh.EncounterService.GetEncounterByID(encounterID)
	if err != nil {
		http.Error(w, "Failed to get encounter by ID", http.StatusInternalServerError)
		return
	}

	// Encode encounter into JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(encounter)
}
