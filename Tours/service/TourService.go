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

// CreateTour creates a new tour
func (ts *TourService) CreateTour(tour *model.Tour) error {
	return ts.TourRepository.CreateTour(tour)
}
