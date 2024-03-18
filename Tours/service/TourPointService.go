package service

import (
	"tours/model"
	"tours/repo"
)

type TourPointService struct {
	TourPointRepository *repo.TourPointRepository
}

func NewTourPointService(tr *repo.TourPointRepository) *TourPointService {
	return &TourPointService{
		TourPointRepository: tr,
	}
}

// CreateTourPoint creates a new tour point
func (tps *TourPointService) CreateTourPoint(tourPoint *model.TourPoint) error {
	return tps.TourPointRepository.CreateTourPoint(tourPoint)
}
