package repo

import (
	"tours/model"

	"gorm.io/gorm"
)

type TourRepository struct {
	DatabaseConnection *gorm.DB
}

func NewTourRepository(db *gorm.DB) *TourRepository {
	return &TourRepository{
		DatabaseConnection: db,
	}
}

func (tr *TourRepository) CreateTour(tour *model.Tour) error {
	return tr.DatabaseConnection.Create(tour).Error
}
