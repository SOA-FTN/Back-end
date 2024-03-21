package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tours/model"
	"tours/service"
)

type TourReviewHandler struct {
	TourReviewService *service.TourReviewService
}

func NewTourReviewHandler(trs *service.TourReviewService) *TourReviewHandler {
	return &TourReviewHandler{
		TourReviewService: trs,
	}
}

func (trh *TourReviewHandler) CreateTourReviewHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the request body into a TourReview object
	var tourReview model.TourReview
	err := json.NewDecoder(r.Body).Decode(&tourReview)
	if err != nil {
		http.Error(w, "Failed to decode request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Call the service method to create the tour review
	err = trh.TourReviewService.CreateTourReview(&tourReview)
	if err != nil {
		http.Error(w, "Failed to create tour review: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Tour review created successfully"})
}

func (trh *TourReviewHandler) GetTourReviewsByTourIDHandler(w http.ResponseWriter, r *http.Request) {
	tourIDStr := r.URL.Query().Get("tourID")
	tourID, err := strconv.Atoi(tourIDStr)
	if err != nil {
		http.Error(w, "Invalid tourID", http.StatusBadRequest)
		return
	}

	tourReviews, err := trh.TourReviewService.GetTourReviewsByTourID(tourID)
	if err != nil {
		http.Error(w, "Failed to retrieve tour reviews: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert tourReviews to JSON
	jsonData, err := json.Marshal(tourReviews)
	if err != nil {
		http.Error(w, "Failed to marshal tour reviews to JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Write JSON data to response body
	_, err = w.Write(jsonData)
	if err != nil {
		http.Error(w, "Failed to write JSON data to response body: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
