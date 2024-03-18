package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"tours/model"
	"tours/service"
)

type TourPointHandler struct {
	TourPointService *service.TourPointService
}

func NewTourPointHandler(tps *service.TourPointService) *TourPointHandler {
	return &TourPointHandler{
		TourPointService: tps,
	}
}

// CreateTourPointHandler handles creating a new tour point
func (tph *TourPointHandler) CreateTourPointHandler(w http.ResponseWriter, r *http.Request) {
	var tourPoint model.TourPoint
	log.Println("aa")
	if err := json.NewDecoder(r.Body).Decode(&tourPoint); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Println(err)
		return
	}

	if err := tph.TourPointService.CreateTourPoint(&tourPoint); err != nil {
		http.Error(w, "Failed to create tour point", http.StatusInternalServerError)
		log.Println(err)

		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tourPoint)
}

func (th *TourPointHandler) GetTourPointsByTourIDHandler(w http.ResponseWriter, r *http.Request) {
	tourIDStr := r.URL.Query().Get("tourId")
	tourID, err := strconv.ParseInt(tourIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid tourID", http.StatusBadRequest)
		return
	}

	tourPoints, err := th.TourPointService.GetTourPointsByTourID(tourID)
	if err != nil {
		http.Error(w, "Failed to get tour points by TourID", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tourPoints)
}
