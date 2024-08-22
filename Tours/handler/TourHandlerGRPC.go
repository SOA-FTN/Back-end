package handler

import (
	"context"
	"fmt"
	"log"
	"tours/model"
	TOURS "tours/proto"
	"tours/service"

	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type TourHandlerGRPC struct {
	TourService *service.TourService
	TourPointService *service.TourPointService
	TOURS.UnimplementedToursServiceServer
}

func NewTourHandlerGRPC(to *service.TourService,tp *service.TourPointService) *TourHandlerGRPC {
	return &TourHandlerGRPC{
		TourService :to,
		TourPointService: tp,
	}
}

func (th *TourHandlerGRPC) TourCreationRpc(ctx context.Context,req *TOURS.TourCreationRequest) (*TOURS.TourCreationResponse,error){
	
	difficultyLevel := model.ToModelDifficultyLevel(req.DifficultyLevel);
	status := model.ToModelStatus(req.Status)

	tour := model.Tour {
		Name:            req.Name,
        DifficultyLevel: difficultyLevel,  // Use the converted value here
        Description:     req.Description,
        Price:           int(req.Price),
        UserId:          int(req.UserId),
		Status:  		 status,	
	}

	TourId,err:= th.TourService.CreateTour(&tour)
	if(err!= nil){
		return nil, fmt.Errorf("Error while registering a new user: %v", err)
	}

	for _,tp := range req.TourPoints {
		tourPoint := model.TourPoint{
            TourId:      int64(TourId), // Postavljanje ID-a kreirane ture
            Name:        tp.Name,
            Description: &tp.Description,
            Latitude:    tp.Latitude,
            Longitude:   tp.Longitude,
            ImageUrl:    tp.ImageUrl,
			Type: model.TourPointType(tp.Type),
        }
		err := th.TourPointService.CreateTourPoint(&tourPoint)
        if err != nil {
            return nil, fmt.Errorf("Error while creating tour point: %v", err)
        }
	}
	return &TOURS.TourCreationResponse{
		Message: "Tour created successfully",
	}, nil
}

func (th *TourHandlerGRPC) GetPublishedToursRpc(ctx context.Context , req *emptypb.Empty) (*TOURS.GetAllToursResponse,error){
	tours, err := th.TourService.GetPublishedTours()
	if err != nil {
        return nil, err
    }
	var protoTours []*TOURS.Tour
	for _, tour := range tours {
		// Pozivanje GetTourPointsByTourID i rukovanje sa greškom
		tourPoints, err := th.TourPointService.GetTourPointsByTourID(int64(tour.ID))
		if err != nil {
			return nil, err
		}

		// Konverzija TourPoints u proto format
		var protoTourPoints []*TOURS.TourPoint
		for _, point := range tourPoints {
			protoTourPoint := &TOURS.TourPoint{
				TourId:      uint64(point.TourId),
				Name:        point.Name,
				Description: *point.Description,
				Latitude:    point.Latitude,
				Longitude:   point.Longitude,
				ImageUrl:    point.ImageUrl,
				Type: TOURS.TourPointType(point.Type),
			}
			protoTourPoints = append(protoTourPoints, protoTourPoint)
		}

		// Kreiranje protoTour objekta
		protoTour := &TOURS.Tour{
			Id:               uint64(tour.ID), // Pretpostavljam da je tour.ID tipa uint64
			Name:             tour.Name,
			DifficultyLevel:  model.ToProtoDifficultyLevel(tour.DifficultyLevel),
			Description:      tour.Description,
			Status:          model.ToProtoStatus(tour.Status),
			Price:            uint64(tour.Price),
			UserId:           uint64(tour.UserId),
			TourPoints:       protoTourPoints,
			
			// PublishedDateTime: model.toProtoTimestamp(tour.PublishedDateTime),
			// ArchivedDateTime:  toProtoTimestamp(tour.ArchivedDateTime),
		}
		protoTours = append(protoTours, protoTour)
	}

	response := &TOURS.GetAllToursResponse{
		Tours: protoTours,
	}
	return response, nil
}

func (th *TourHandlerGRPC) PurchaseToursRpc(ctx context.Context,req *TOURS.PurchaseToursRequest) (*TOURS.PurchaseToursResponse,error){
	userId := int(req.UserId)
	tourIds := make([]int, len(req.TourId))


	for i, id := range req.TourId {
		tourIds[i] = int(id)
	}
	err := th.TourService.PurchaseTours(userId, tourIds)
	if err != nil {
		log.Printf("Failed to purchase tours for user %d: %v", userId, err)
		return &TOURS.PurchaseToursResponse{
			Message: "Failed to purchase tours",
		}, err
	}

	return &TOURS.PurchaseToursResponse{
		Message: "Tours purchased successfully",
	}, nil
}

func (th *TourHandlerGRPC) GetPurchasedToursRpc(ctx context.Context, req *TOURS.UserPurchaseTourRequest) (*TOURS.UserPurchaseTourResponse, error) {
	userID := int(req.UserId)

	tours, err := th.TourService.GetPurchasedToursByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("Error getting purchased tours: %v", err)
	}

	var protoTours []*TOURS.Tour
	for _, tour := range tours {
		// Pozivanje GetTourPointsByTourID i rukovanje sa greškom
		tourPoints, err := th.TourPointService.GetTourPointsByTourID(int64(tour.ID))
		if err != nil {
			return nil, err
		}

		// Konverzija TourPoints u proto format
		var protoTourPoints []*TOURS.TourPoint
		for _, point := range tourPoints {
			protoTourPoint := &TOURS.TourPoint{
				TourId:      uint64(point.TourId),
				Name:        point.Name,
				Description: *point.Description,
				Latitude:    point.Latitude,
				Longitude:   point.Longitude,
				ImageUrl:    point.ImageUrl,
				Type: TOURS.TourPointType(point.Type),
			}
			protoTourPoints = append(protoTourPoints, protoTourPoint)
		}

		// Kreiranje protoTour objekta
		protoTour := &TOURS.Tour{
			Id:               uint64(tour.ID), // Pretpostavljam da je tour.ID tipa uint64
			Name:             tour.Name,
			DifficultyLevel:  model.ToProtoDifficultyLevel(tour.DifficultyLevel),
			Description:      tour.Description,
			Status:          model.ToProtoStatus(tour.Status),
			Price:            uint64(tour.Price),
			UserId:           uint64(tour.UserId),
			TourPoints:       protoTourPoints,
			
			// PublishedDateTime: model.toProtoTimestamp(tour.PublishedDateTime),
			// ArchivedDateTime:  toProtoTimestamp(tour.ArchivedDateTime),
		}
		protoTours = append(protoTours, protoTour)
	}

	response := &TOURS.UserPurchaseTourResponse{
		Tours: protoTours,
	}
	return response, nil
}

func (th *TourHandlerGRPC) GetAuthorToursRpc(ctx context.Context, req *TOURS.AuthorToursRequest) (*TOURS.AuthorToursResponse, error) {
    authorId := req.UserId
    tours, err := th.TourService.GetToursByUserID(int(authorId))
    if err != nil {
        return nil, fmt.Errorf("Error getting author tours: %v", err)
    }

    var protoTours []*TOURS.UpdateTourRequest
    for _, tour := range tours {
        // Pozivanje GetTourPointsByTourID i rukovanje sa greškom
        tourPoints, err := th.TourPointService.GetTourPointsByTourID(int64(tour.ID))
        if err != nil {
            return nil, err
        }

        // Konverzija TourPoints u proto format
        var protoTourPoints []*TOURS.UpdateTourPoint
        for _, point := range tourPoints {
            protoTourPoint := &TOURS.UpdateTourPoint{
                Id:          uint64(point.ID), // Dodajte ID za UpdateTourPoint
                TourId:      uint64(point.TourId),
                Name:        point.Name,
                Description: *point.Description,
                Latitude:    point.Latitude,
                Longitude:   point.Longitude,
                ImageUrl:    point.ImageUrl,
                Type:        TOURS.TourPointType(point.Type),
            }
            protoTourPoints = append(protoTourPoints, protoTourPoint)
        }

        // Kreiranje UpdateTourRequest objekta umesto Tour objekta
        protoTour := &TOURS.UpdateTourRequest{
            Id:              uint64(tour.ID),
            Name:            tour.Name,
            DifficultyLevel: TOURS.DifficultyLevel(tour.DifficultyLevel),
            Description:     tour.Description,
            Status:          TOURS.TourStatus(tour.Status),
            Price:           uint64(tour.Price),
            UserId:          uint64(tour.UserId),
            TourPoints:      protoTourPoints,
            // PublishedDateTime: tour.PublishedDateTime,
            // ArchivedDateTime:  tour.ArchivedDateTime,
        }
        protoTours = append(protoTours, protoTour)
    }

    response := &TOURS.AuthorToursResponse{
        Tours: protoTours,
    }
    return response, nil
}

func (th *TourHandlerGRPC) UpdateTourRpc(ctx context.Context, req *TOURS.UpdateTourRequest) (*TOURS.TourCreationResponse, error) {
	
	// Pretvaranje enuma u model vrednosti
	difficultyLevel := model.ToModelDifficultyLevel(req.DifficultyLevel)
	status := model.ToModelStatus(req.Status)

	// Kreiranje osnovne Tour strukture
	tour := model.Tour{
		Model: gorm.Model{
			ID: uint(req.Id), // Pretpostavljam da `req.Id` sadrži ID iz Protobuf-a
		},
		Name:            req.Name,
		DifficultyLevel: difficultyLevel,
		Description:     req.Description,
		Price:           int(req.Price),
		UserId:          int(req.UserId),
		Status:          status,
	}

	// Ažuriranje ture
	updatedTour, err := th.TourService.UpdateTour(&tour)
	if err != nil {
		return nil, fmt.Errorf("Error while updating the tour: %v", err)
	}

	// Konvertovanje TourPoints iz req u model.TourPoint
	var tourPoints []model.TourPoint
	for _, tp := range req.TourPoints {
		tourPoints = append(tourPoints, model.TourPoint{
			Model:       gorm.Model{ID: uint(tp.Id)}, // Ako ID dolazi iz req, postavi ga ovde
			TourId:      int64(updatedTour.ID),
			Name:        tp.Name,
			Description: &tp.Description,
			Latitude:    tp.Latitude,
			Longitude:   tp.Longitude,
			ImageUrl:    tp.ImageUrl,
			Type: model.TourPointType(tp.Type), // Pretvaranje tipa
		})
	}

	for _, tp := range tourPoints {
		fmt.Printf("TourPoint ID: %d, Name: %s, Description: %s, Latitude: %f, Longitude: %f, Type: %d\n",
			tp.ID, tp.Name, *tp.Description, tp.Latitude, tp.Longitude,  tp.Type)
	}

	// Ažuriranje TourPointova za turu
	_, err = th.TourPointService.UpdateTourPointsOfTour(int64(updatedTour.ID), tourPoints)
	if err != nil {
		return nil, fmt.Errorf("Error while updating tour points: %v", err)
	}

	// Vraćanje odgovora
	return &TOURS.TourCreationResponse{
		TourId: uint64(updatedTour.ID),
	}, nil
}


