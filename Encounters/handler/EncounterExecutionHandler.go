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
	"strconv"

	"github.com/gorilla/mux"
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

	var requestBody bytes.Buffer

	if _, err := io.Copy(&requestBody, r.Body); err != nil {
		log.Println("Failed to read request body:", err)
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	r.Body = io.NopCloser(&requestBody)

	var enc model.EncounterExecution
	if err := json.NewDecoder(&requestBody).Decode(&enc); err != nil {
		log.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := eh.EncounterExecutionService.CreateEncounterExecution(&enc); err != nil {
		http.Error(w, "Failed to create encounter", http.StatusInternalServerError)

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

func (eh *EncounterExecutionHandler) GetEncounterExecutionByUserIDAndNotCompletedHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	userIDStr, ok := vars["userId"]
	if !ok {
		http.Error(w, "User ID not provided", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid userID", http.StatusBadRequest)
		return
	}

	encounter, err := eh.EncounterExecutionService.GetEncounterExecutionByUserIDAndNotCompleted(userID)
	if err != nil {
		http.Error(w, "Failed to get encounter by UserID", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(encounter)
}

func (eeh *EncounterExecutionHandler) UpdateEncounterExecutionHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	userIDStr, ok := vars["userId"]
	if !ok {
		http.Error(w, "User ID not provided", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid userID", http.StatusBadRequest)
		return
	}

	if err := eeh.EncounterExecutionService.UpdateEncounterExecution(userID); err != nil {
		http.Error(w, "Failed to update encounter execution", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
