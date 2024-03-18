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

	newTour := model.Tour{
		Name:            tour.Name,
		DifficultyLevel: tour.DifficultyLevel,
		Description:     tour.Description,
		Status:          model.Draft,
		Price:           tour.Price,
		UserId:          tour.UserId,
	}
	err := ts.TourRepository.CreateTour(&newTour)
	if err != nil {
		return err
	}
	return nil
}

func (ts *TourService) GetToursByUserID(userID int) ([]model.Tour, error) {
	// Call repository function to get tours by userID
	tours, err := ts.TourRepository.GetToursByUserID(userID)
	if err != nil {
		return nil, err
	}

	return tours, nil
}

func ConvertDifficultyLevelToInt(difficultyLevel string) int {
	switch difficultyLevel {
	case "Easy":
		return 0
	case "Moderate":
		return 1
	case "Difficult":
		return 2
	default:
		return -1
	}
}

// ConvertStatusToInt converts the string status to an integer.
// It returns -1 if the status is not recognized.
func ConvertStatusToInt(status string) int {
	switch status {
	case "Draft":
		return 0
	case "Published":
		return 1
	case "Archived":
		return 2
	default:
		return -1
	}
}
