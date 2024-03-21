package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"tours/model"
	"tours/service"

	"github.com/gorilla/mux"
)

type TourHandler struct {
	TourService *service.TourService
}

func NewTourHandler(ts *service.TourService) *TourHandler {
	return &TourHandler{
		TourService: ts,
	}
}

func (th *TourHandler) CreateTourHandler(w http.ResponseWriter, r *http.Request) {

	var req model.CreateTourRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tour := model.Tour{
		Name:              req.Name,
		Description:       req.Description,
		DifficultyLevel:   model.DifficultyLevel(service.ConvertDifficultyLevelToInt(req.DifficultyLevel)),
		TStatus:           model.TourStatus(service.ConvertStatusToInt(req.Status)),
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

func (th *TourHandler) GetPublishedToursHandler(w http.ResponseWriter, r *http.Request) {
	// Call the service method to get approved tours
	tours, err := th.TourService.GetPublishedTours()
	if err != nil {
		http.Error(w, "Failed to retrieve approved tours: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the retrieved tours as JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tours)
}

func (th *TourHandler) GetToursByUserIDHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("userId")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid userID", http.StatusBadRequest)
		return
	}

	tours, err := th.TourService.GetToursByUserID(userID)
	if err != nil {
		http.Error(w, "Failed to get tours by UserID", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tours)
}

func (th *TourHandler) UpdateTourHandler(w http.ResponseWriter, r *http.Request) {
	var tour model.Tour
	err := json.NewDecoder(r.Body).Decode(&tour)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedTour, err := th.TourService.UpdateTour(&tour)
	if err != nil {
		http.Error(w, "Failed to update tour", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTour)
}

func (th *TourHandler) PublishTourHandler(w http.ResponseWriter, r *http.Request) {
	// Extract tourID from request URL parameter
	vars := mux.Vars(r)
	tourIDStr, ok := vars["tourID"]
	if !ok {
		http.Error(w, "tourID not found in URL path", http.StatusBadRequest)
		return
	}

	tourID, err := strconv.ParseInt(tourIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid tourID", http.StatusBadRequest)
		return
	}

	// Call the service method to publish tour
	err = th.TourService.PublishTour(tourID)
	if err != nil {
		http.Error(w, "Failed to publish tour: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Tour published successfully"})
}

func (th *TourHandler) ArchiveTourHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tourIDStr, ok := vars["tourID"]
	if !ok {
		http.Error(w, "tourID not found in URL path", http.StatusBadRequest)
		return
	}

	tourID, err := strconv.ParseInt(tourIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid tourID", http.StatusBadRequest)
		return
	}

	err = th.TourService.ArchiveTour(tourID)
	if err != nil {
		http.Error(w, "Failed to archive tour: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Tour archived successfully"})
}
