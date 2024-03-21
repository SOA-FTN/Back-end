package repo

import (
	"tours/model"

	"gorm.io/gorm"
)

type TourReviewRepository struct {
	DatabaseConnection *gorm.DB
}

func NewTourReviewRepository(db *gorm.DB) *TourReviewRepository {
	return &TourReviewRepository{
		DatabaseConnection: db,
	}
}

func (tr *TourReviewRepository) CreateTourReview(review *model.TourReview) error {
	dbResult := tr.DatabaseConnection.Create(review)
	if dbResult.Error != nil {
		return dbResult.Error
	}
	return nil
}

func (tr *TourReviewRepository) GetAllTourReviews() ([]*model.TourReview, error) {
	var tourReviews []*model.TourReview
	if err := tr.DatabaseConnection.Find(&tourReviews).Error; err != nil {
		return nil, err
	}
	return tourReviews, nil
}

func (trr *TourReviewRepository) GetTourReviewByTourIDAndUserID(tourID, userID int) (*model.TourReview, error) {
	var tourReview model.TourReview
	if err := trr.DatabaseConnection.Where("tour_id = ? AND user_id = ?", tourID, userID).First(&tourReview).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &tourReview, nil
}

func (trr *TourReviewRepository) GetTourReviewsByTourID(tourID int) ([]model.TourReview, error) {
	var tourReviews []model.TourReview
	if err := trr.DatabaseConnection.Where("tour_id = ?", tourID).Find(&tourReviews).Error; err != nil {
		return nil, err
	}
	return tourReviews, nil
}
