package model

import "time"

// TourReview represents a review for a tour.
type TourReview struct {
	Rating      int       // Rating given for the tour
	Comment     string    // Comment left by the user
	UserID      int       // ID of the user who left the review
	TourID      int       // ID of the tour being reviewed
	ImageURL    string    // URL of an image related to the review (optional)
	VisitDate   time.Time // Date when the user visited the tour
	CommentDate time.Time // Date when the user left the comment
}
