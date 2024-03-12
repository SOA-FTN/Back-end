package handler

import (
	"encoding/json"
	"net/http"
	"tours/model"
	"tours/service"
)

type TourHandler struct {
	TourService *service.TourService
}

func NewTourHandler(ts *service.TourService) *TourHandler {
	return &TourHandler{
		TourService: ts,
	}
}

// CreateTourHandler handles creating a new tour
func (th *TourHandler) CreateTourHandler(w http.ResponseWriter, r *http.Request) {
	var tour model.Tour
	if err := json.NewDecoder(r.Body).Decode(&tour); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := th.TourService.CreateTour(&tour); err != nil {
		http.Error(w, "Failed to create tour", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tour)
}
