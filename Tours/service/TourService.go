package service

import (
	"tours/model"
	"tours/repo"
)

type TourService struct {
	TourRepository *repo.TourRepository
}

func NewTourService(tr *repo.TourRepository) *TourService {
	return &TourService{
		TourRepository: tr,
	}
}

func (ts *TourService) CreateTour(tour *model.Tour) error {
	return ts.TourRepository.CreateTour(tour)
}

func (ts *TourService) GetToursByUserID(userID int) ([]model.Tour, error) {
	// Call repository function to get tours by userID
	tours, err := ts.TourRepository.GetToursByUserID(userID)
	if err != nil {
		return nil, err
	}

	return tours, nil
}
