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

func (ts *TourPointService) UpdateTourPointsOfTour(tourId int64, tourPoints []model.TourPoint) ([]model.TourPoint, error) {
	// Dohvatanje trenutnih tour pointova iz baze
	currentTourPoints, err := ts.GetTourPointsByTourID(tourId)
    if err != nil {
        return nil, err
    }

	// Kreiranje mapa za jednostavnije poređenje
	currentTourPointsMap := make(map[int64]model.TourPoint)
    newTourPointsMap := make(map[int64]model.TourPoint)

	for _, tp := range currentTourPoints {
        currentTourPointsMap[int64(tp.ID)] = tp
    }

    // Mapiranje novih tour pointova
    for _, newTP := range tourPoints {
        newTourPointsMap[int64(newTP.ID)] = newTP
    }

	// Iteriranje kroz nove tour pointove
	for id, currentTP := range currentTourPointsMap {
        if newTP, found := newTourPointsMap[id]; found {
            // Ako tour point postoji u obe liste, ažuriraj ga
            currentTP.Name = newTP.Name
            currentTP.Description = newTP.Description
            currentTP.Latitude = newTP.Latitude
            currentTP.Longitude = newTP.Longitude
            currentTP.ImageUrl = newTP.ImageUrl
            currentTP.Type = newTP.Type

            if err := ts.TourPointRepository.UpdateTourPoint(&currentTP); err != nil {
                return nil, err
            }
        } else {
            // Ako tour point ne postoji u novoj listi, obriši ga
            if err := ts.TourPointRepository.DeleteTourPoint(int64(currentTP.ID)); err != nil {
                return nil, err
            }
        }
    }

    // Iteriranje kroz nove tour pointove da se vidi koji treba da se dodaju
    for _, newTP := range tourPoints {
        if _, found := currentTourPointsMap[int64(newTP.ID)]; !found {
            newTP.TourId = tourId
            if err := ts.TourPointRepository.CreateTourPoint(&newTP); err != nil {
                return nil, err
            }
        }
    }

    // Vraćanje ažurirane liste tour pointova
    updatedTourPoints, err := ts.GetTourPointsByTourID(tourId)
    if err != nil {
        return nil, err
    }

    return updatedTourPoints, nil
}
