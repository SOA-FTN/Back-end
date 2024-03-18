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

func (tps *TourPointService) CreateTourPoint(tourPoint *model.TourPoint) error {
	return tps.TourPointRepository.CreateTourPoint(tourPoint)
}

func (ts *TourPointService) GetTourPointsByTourID(tourID int64) ([]model.TourPoint, error) {
	tourPoints, err := ts.TourPointRepository.GetTourPointsByTourID(tourID)
	if err != nil {
		return nil, err
	}
	return tourPoints, nil
}
