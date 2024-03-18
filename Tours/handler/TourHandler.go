package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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

	var req model.CreateTourRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Map fields to the Tour struct
	tour := model.Tour{
		Name:              req.Name,
		Description:       req.Description,
		DifficultyLevel:   model.DifficultyLevel(service.ConvertDifficultyLevelToInt(req.DifficultyLevel)),
		Status:            model.TourStatus(service.ConvertStatusToInt(req.Status)),
		Price:             req.Price,
		UserId:            req.UserID,
		ArchivedDateTime:  req.ArchivedDateTime,
		PublishedDateTime: req.PublishedDateTime,
	}

	if err := th.TourService.CreateTour(&tour); err != nil {
		http.Error(w, "Failed to create tour", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tour)
}

// Handler Function to Get Tours by UserID
func (th *TourHandler) GetToursByUserIDHandler(w http.ResponseWriter, r *http.Request) {
	// Extract userID from request, assuming it's a query parameter
	userIDStr := r.URL.Query().Get("userId")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid userID", http.StatusBadRequest)
		return
	}

	// Call service function to get tours by userID
	tours, err := th.TourService.GetToursByUserID(userID)
	if err != nil {
		http.Error(w, "Failed to get tours by UserID", http.StatusInternalServerError)
		return
	}

	// Encode tours into JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tours)
}
