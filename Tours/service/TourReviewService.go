package service

import (
	"errors"
	"tours/model"
	"tours/repo"
)

type TourReviewService struct {
	TourReviewRepository *repo.TourReviewRepository
}

func NewTourReviewService(tr *repo.TourReviewRepository) *TourReviewService {
	return &TourReviewService{
		TourReviewRepository: tr,
	}
}

func (trs *TourReviewService) CreateTourReview(review *model.TourReview) error {
	// Check if there's already a TourReview with the same combination of tourId and userId
	existingReview, err := trs.TourReviewRepository.GetTourReviewByTourIDAndUserID(review.TourID, review.UserID)
	if err != nil {
		return err
	}
	// If an existing review with the same tourId and userId exists, return an error
	if existingReview != nil {
		return errors.New("another review already exists for the same tour and user")
	}

	// Create the TourReview if no duplicate is found
	err = trs.TourReviewRepository.CreateTourReview(review)
	if err != nil {
		return err
	}
	return nil
}

func (trs *TourReviewService) GetTourReviewsByTourID(tourID int) ([]model.TourReview, error) {
	return trs.TourReviewRepository.GetTourReviewsByTourID(tourID)
}
