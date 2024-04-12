package handler

import (
	"context"
	"encounters/model"
	"encounters/service"

	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type KeyProduct struct{}

type EncounterHandler struct {
	logger *log.Logger
	EncounterService *service.EncounterService
}

func NewEncounterHandler(l *log.Logger, es  *service.EncounterService) *EncounterHandler {
	return &EncounterHandler{
		l,es,
	}
}
/*
func (eh *EncounterHandler) CreateEncounterHandler(w http.ResponseWriter, r *http.Request) {

	var enc model.CreateEncounter
	if err := json.NewDecoder(r.Body).Decode(&enc); err != nil {
		log.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

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
*/

func (eh *EncounterHandler) CreateEncounterHandler(rw http.ResponseWriter, h *http.Request) {
	encounter := h.Context().Value(KeyProduct{}).(*model.CreateEncounter)
	log.Println(encounter)
	if err := eh.EncounterService.CreateEncounter(encounter); err != nil {
		http.Error(rw, "Failed to create encounter", http.StatusInternalServerError)
		log.Println("ne")
		return
	}
	rw.WriteHeader(http.StatusCreated)
	//json.NewEncoder(h).Encode(enc)
}

func (eh *EncounterHandler) GetAllEncountersHandler(w http.ResponseWriter, r *http.Request) {
	encounters, err := eh.EncounterService.GetAllEncounters()
	if err != nil {
		http.Error(w, "Failed to get encounters", http.StatusInternalServerError)
		eh.logger.Println("Failed to get encounters:", err)
		return
	}
	/*
	// Convert encounters to JSON and send response
	response, err := json.Marshal(encounters)
	if err != nil {
		http.Error(w, "Failed to marshal encounters", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	*/
	if encounters == nil {
		http.Error(w, "No encounters found", http.StatusNotFound)
		eh.logger.Println("No encounters found")
		return
	}

	err = encounters.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to convert to json", http.StatusInternalServerError)
		eh.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (eh *EncounterHandler) GetEncounterByIDHandler(rw http.ResponseWriter, h *http.Request) {

	vars := mux.Vars(h)
	encounterid := vars["encounterId"]
	/*
	encounterID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(rw, "Invalid encounter ID", http.StatusBadRequest)
		return
	}
	*/
	encounter , err := eh.EncounterService.GetEncounterByID(encounterid)
	if err != nil {
		eh.logger.Print("Database exception: ", err)
	}

	if encounter == nil {
		http.Error(rw, "Encounter with given id not found", http.StatusNotFound)
		eh.logger.Printf("Encounter with id: '%s' not found", encounterid)
		return
	}

	err =encounter.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		eh.logger.Fatal("Unable to convert to json :", err)
		return
	}

	/*
	vars := mux.Vars(r)
	encounterIDStr, ok := vars["encounterId"]
	if !ok {
		http.Error(w, "Encounter ID not provided", http.StatusBadRequest)
		return
	}

	encounterID, err := strconv.Atoi(encounterIDStr)
	if err != nil {
		http.Error(w, "Invalid encounter ID", http.StatusBadRequest)
		return
	}

	encounter, err := eh.EncounterService.GetEncounterByID(encounterID)
	if err != nil {
		http.Error(w, "Failed to get encounter by ID", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(encounter)
	*/
}


func(e *EncounterHandler) MIddlewareEncounterDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		log.Println(h.Body)
		log.Println("==============")
		encounter := &model.CreateEncounter{}
		err:=encounter.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			e.logger.Fatal(err)
			return
		}
		log.Println(encounter)
		log.Println("=============")
		ctx := context.WithValue(h.Context(), KeyProduct{}, encounter)
		h = h.WithContext(ctx)

		next.ServeHTTP(rw, h)
	})
}


func (e *EncounterHandler) MiddlewareContentTypeSet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		e.logger.Println("Method [", h.Method, "] - Hit path :", h.URL.Path)

		rw.Header().Add("Content-Type", "application/json")

		next.ServeHTTP(rw, h)
	})
}
