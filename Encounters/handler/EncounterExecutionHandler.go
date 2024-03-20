package handler

import (
	"encoding/json"
	"encounters/service"
	"net/http"
)

type EncounterExecutionHandler struct {
	EncounterExecutionService *service.EncounterExecutionService
}

func NewEncounterExecutionHandler(es *service.EncounterExecutionService) *EncounterExecutionHandler {
	return &EncounterExecutionHandler{
		EncounterExecutionService: es,
	}
}

func (eh *EncounterExecutionHandler) GetAllEncounterExecutionsHandler(w http.ResponseWriter, r *http.Request) {
	encounters, err := eh.EncounterExecutionService.GetAllEncounterExecutions()
	if err != nil {
		http.Error(w, "Failed to get encounters", http.StatusInternalServerError)
		return
	}

	// Convert encounters to JSON and send response
	response, err := json.Marshal(encounters)
	if err != nil {
		http.Error(w, "Failed to marshal encounters", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
