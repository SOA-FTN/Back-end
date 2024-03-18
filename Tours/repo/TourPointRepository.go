package repo

import (
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
