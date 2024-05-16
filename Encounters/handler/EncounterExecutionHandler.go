package handler

import (
	"context"
	"encounters/model"
	"encounters/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type EncounterExecutionHandler struct {
	logger *log.Logger
	EncounterExecutionService *service.EncounterExecutionService
}

func NewEncounterExecutionHandler(l *log.Logger,es *service.EncounterExecutionService) *EncounterExecutionHandler {
	return &EncounterExecutionHandler{
		l,es,
	}
}

func (eh *EncounterExecutionHandler) CreateEncounterExecutionHandler(rw http.ResponseWriter, h *http.Request) {
	encounterExecution := h.Context().Value(KeyProduct{}).(*model.EncounterExecution)
	log.Println("JEBEM TI SE SA MAMAROM")
	log.Println(encounterExecution)
	if err := eh.EncounterExecutionService.CreateEncounterExecution(encounterExecution); err != nil {
		http.Error(rw, "Failed to create encounter execution", http.StatusInternalServerError)
		log.Println("ne")
	}

	rw.WriteHeader(http.StatusCreated)
	
}

func (eh *EncounterExecutionHandler) GetAllEncounterExecutionsHandler(w http.ResponseWriter, r *http.Request) {
	executions, err := eh.EncounterExecutionService.GetAllEncounterExecutions()
	if err != nil {
		http.Error(w, "Failed to get encounter executions", http.StatusInternalServerError)
		eh.logger.Println("Failed to get encounters executions:", err)
		return
	}

	if executions == nil {
		http.Error(w, "No encounter executions found", http.StatusNotFound)
		eh.logger.Println("No encounter executions  found")
		return
	}

	err = executions.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to convert to json", http.StatusInternalServerError)
		eh.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (eh *EncounterExecutionHandler) GetEncounterExecutionByUserIDAndNotCompletedHandler(rw http.ResponseWriter, h *http.Request) {

	vars := mux.Vars(h)
	userIDStr, ok := vars["userId"]
	if !ok {
		http.Error(rw, "User ID not provided", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(rw, "Invalid userID", http.StatusBadRequest)
		return
	}

	encounter, err := eh.EncounterExecutionService.GetEncounterExecutionByUserIDAndNotCompleted(userID)
	if err != nil {
		http.Error(rw, "Failed to get encounter by UserID", http.StatusInternalServerError)
		return
	}

	if encounter == nil {
		http.Error(rw, "Encounter with given user id not found", http.StatusNotFound)
		eh.logger.Printf("Encounter with userid: '%s' not found", userID)
		return
	}

	err =encounter.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		eh.logger.Fatal("Unable to convert to json :", err)
		return
	}
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


func(e *EncounterExecutionHandler) MIddlewareEncounterExecutionDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		log.Println("STA SE DESAVA BAJO")
		log.Println(h.Body)
		log.Println("==============")
		encounterExecution := &model.EncounterExecution{}
		err:=encounterExecution.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			e.logger.Fatal(err)
			return
		}
		log.Println("KOJI KURAC DRUZEEE")
		log.Println(encounterExecution)
		log.Println("=============")
		ctx := context.WithValue(h.Context(), KeyProduct{}, encounterExecution)
		h = h.WithContext(ctx)

		next.ServeHTTP(rw, h)
	})
}