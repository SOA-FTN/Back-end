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
	dbResult := tr.DatabaseConnection.Create(tour)
	if dbResult.Error != nil {
		return dbResult.Error
	}
	return nil
}

func (tr *TourRepository) GetToursByUserID(userID int) ([]model.Tour, error) {
	var tours []model.Tour
	if err := tr.DatabaseConnection.Where("user_id = ?", userID).Find(&tours).Error; err != nil {
		return nil, err
	}
	return tours, nil
}

func (tr *TourRepository) UpdateTour(tour *model.Tour) (*model.Tour, error) {
	dbResult := tr.DatabaseConnection.Model(&model.Tour{}).Where("id = ?", tour.ID).Updates(tour)
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	return tour, nil
}

func (tr *TourRepository) FindTourByID(id int) (*model.Tour, error) {
	var tour model.Tour
	if err := tr.DatabaseConnection.First(&tour, id).Error; err != nil {
		return nil, err
	}
	return &tour, nil
}

func (tr *TourRepository) GetPublishedTours() ([]model.Tour, error) {
	var tours []model.Tour
	err := tr.DatabaseConnection.Where("t_status = ?", 1).Find(&tours).Error
	if err != nil {
		return nil, err
	}
	return tours, nil
}
