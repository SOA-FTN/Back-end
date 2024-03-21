package service

import (
	"errors"
	"log"
	"tours/model"
	"tours/repo"
)

type TourService struct {
	TourRepository      *repo.TourRepository
	TourPointRepository *repo.TourPointRepository
}

func NewTourService(tr *repo.TourRepository, tps *repo.TourPointRepository) *TourService {
	return &TourService{
		TourRepository:      tr,
		TourPointRepository: tps,
	}
}

func (ts *TourService) GetPublishedTours() ([]model.Tour, error) {
	tours, err := ts.TourRepository.GetPublishedTours()
	if err != nil {
		return nil, err
	}
	return tours, nil
}

func (ts *TourService) CreateTour(tour *model.Tour) error {

	newTour := model.Tour{
		Name:              tour.Name,
		DifficultyLevel:   tour.DifficultyLevel,
		Description:       tour.Description,
		TStatus:           model.Draft,
		Price:             tour.Price,
		UserId:            tour.UserId,
		ArchivedDateTime:  tour.ArchivedDateTime,
		PublishedDateTime: tour.PublishedDateTime,
	}
	err := ts.TourRepository.CreateTour(&newTour)
	if err != nil {
		return err
	}
	return nil
}

func (ts *TourService) GetToursByUserID(userID int) ([]model.Tour, error) {
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

func ConvertIntToDifficultyLevel(level int) string {
	switch level {
	case 0:
		return "Easy"
	case 1:
		return "Moderate"
	case 2:
		return "Difficult"
	default:
		return "Unknown"
	}
}

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

func (ts *TourService) UpdateTour(tour *model.Tour) (*model.Tour, error) {
	updatedTour, err := ts.TourRepository.UpdateTour(tour)
	if err != nil {
		return nil, err
	}
	return updatedTour, nil
}

func (ts *TourService) PublishTour(tourID int64) error {
	tourPoints, err := ts.TourPointRepository.GetTourPointsByTourID(tourID)

	if err != nil {
		log.Println("Error getting tour points:", err)
		return err
	}

	if len(tourPoints) >= 2 {
		tour, err := ts.TourRepository.FindTourByID(int(tourID))
		if err != nil {
			log.Println("Error finding tour by ID:", err)
			return err
		}

		if tour.Name == "" || tour.Description == "" {
			err := errors.New("tour name or description is empty")
			log.Println("Error: tour name or description is empty")
			return err
		}

		tour.TStatus = 1

		tour, err = ts.TourRepository.UpdateTour(tour)
		if err != nil {
			log.Println("Error updating tour:", err)
			return err
		}
	}

	return nil
}

func (ts *TourService) ArchiveTour(tourID int64) error {

	tour, err := ts.TourRepository.FindTourByID(int(tourID))
	if err != nil {
		log.Println("Error finding tour by ID:", err)
		return err
	}

	tour.TStatus = 2

	tour, err = ts.TourRepository.UpdateTour(tour)
	if err != nil {
		log.Println("Error updating tour:", err)
		return err
	}

	return nil
}
