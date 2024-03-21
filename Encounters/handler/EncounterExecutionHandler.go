package handler

import (
	"bytes"
	"encoding/json"
	"encounters/model"
	"encounters/service"
	"fmt"
	"io"
	"log"
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

func (eh *EncounterExecutionHandler) CreateEncounterExecutionHandler(w http.ResponseWriter, r *http.Request) {
	// Declare a buffer to store the request body
	var requestBody bytes.Buffer
	// Copy the request body into the buffer
	if _, err := io.Copy(&requestBody, r.Body); err != nil {
		log.Println("Failed to read request body:", err)
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	// Log the request body
	log.Println("Request Body:", requestBody.String())

	// Reset the request body so it can be read again later
	r.Body = io.NopCloser(&requestBody)

	var enc model.EncounterExecution
	if err := json.NewDecoder(&requestBody).Decode(&enc); err != nil {
		log.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Println("dosao")
	if err := eh.EncounterExecutionService.CreateEncounterExecution(&enc); err != nil {
		http.Error(w, "Failed to create encounter", http.StatusInternalServerError)
		log.Println("ne")
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(enc)
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

	fmt.Println("Received JSON from front-end:", string(response))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
