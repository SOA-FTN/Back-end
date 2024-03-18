package repo

import (
	"errors"
	"tours/model"

	"gorm.io/gorm"
)

type TourPointRepository struct {
	DatabaseConnection *gorm.DB
}

func NewTourPointRepository(db *gorm.DB) *TourPointRepository {
	return &TourPointRepository{
		DatabaseConnection: db,
	}
}

func (tpr *TourPointRepository) CreateTourPoint(tourPoint *model.TourPoint) error {
	return tpr.DatabaseConnection.Create(tourPoint).Error
}

func (tr *TourPointRepository) GetTourPointsByTourID(tourID int64) ([]model.TourPoint, error) {
	var tourPoints []model.TourPoint
	if err := tr.DatabaseConnection.Where("tour_id = ?", tourID).Find(&tourPoints).Error; err != nil {
		return nil, err
	}
	if len(tourPoints) == 0 {
		return nil, errors.New("no tour points found for the given tour ID")
	}
	return tourPoints, nil
}
