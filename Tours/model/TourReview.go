package model

import "time"

// TourReview represents a review for a tour.
type TourReview struct {
	Rating      int       `json:"rating"`      // Rating given for the tour
	Comment     string    `json:"comment"`     // Comment left by the user
	UserID      int       `json:"userId"`      // ID of the user who left the review
	TourID      int       `json:"tourId"`      // ID of the tour being reviewed
	ImageURL    string    `json:"imageURL"`    // URL of an image related to the review (optional)
	VisitDate   time.Time `json:"visitDate"`   // Date when the user visited the tour
	CommentDate time.Time `json:"commentDate"` // Date when the user left the comment
}
