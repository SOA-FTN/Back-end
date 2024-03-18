package handler

import (
	"encoding/json"
	"net/http"
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
	if err := json.NewDecoder(r.Body).Decode(&tourPoint); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := tph.TourPointService.CreateTourPoint(&tourPoint); err != nil {
		http.Error(w, "Failed to create tour point", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tourPoint)
}
