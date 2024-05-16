package handler

import (
	"context"
	"encounters/model"
	Encounters "encounters/proto"
	"encounters/service"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EncounterHandlergRPC struct {
	logger           *log.Logger
	EncounterService *service.EncounterService
	Encounters.UnimplementedEncounterServiceServer
}

func NewEncounterHandlergRPC(logger *log.Logger, encounterService *service.EncounterService) *EncounterHandlergRPC {
	return &EncounterHandlergRPC{
		logger:           logger,
		EncounterService: encounterService,
	}
}

func (eh *EncounterHandlergRPC) CreateEncounterRpc(ctx context.Context, req *Encounters.CreateEncounterRequest) (*Encounters.CreateEncounterResponse, error) {

	log.Printf("Request received: %+v\n", req)
	encounter := model.CreateEncounter{
		Name:             req.Name,
		Description:      req.Description,
		XpPoints:         req.XpPoints,
		Status:           req.Status.String(),
		Type:             req.Type.String(),
		Latitude:         req.Latitude,
		Longitude:        req.Longitude,
		ShouldBeApproved: req.ShouldBeApproved,
	}
	if eh.EncounterService == nil {
		log.Println("SERVIS  ERRORCINA")
        
    }
	err := eh.EncounterService.CreateEncounter(&encounter)
	if err != nil {
		return nil, err
	}

	return &Encounters.CreateEncounterResponse{}, nil
}

func (eh *EncounterHandlergRPC) GetAllEncountersRpc(ctx context.Context, req *Encounters.GetAllEncountersRequest) (*Encounters.GetAllEncountersResponse, error) {
    encounters, err := eh.EncounterService.GetAllEncounters()
    if err != nil {
        return nil, err
    }

    var protoEncounters []*Encounters.Encounter
    for _, encounter := range encounters {
        protoEncounter := &Encounters.Encounter{
            Id:          encounter.ID.Hex(),
            Name:        encounter.Name,
            Description: encounter.Description,
			XpPoints:    encounter.XpPoints,
			Latitude:    encounter.Latitude,
			Longitude:   encounter.Longitude,
			ShouldBeApproved: encounter.ShouldBeApproved,
			Status:      Encounters.EncounterStatus(encounter.Status),
			Type:        Encounters.EncounterType(encounter.Type),
            
        }
        protoEncounters = append(protoEncounters, protoEncounter)
    }

    response := &Encounters.GetAllEncountersResponse{
        Encounters: protoEncounters,
    }

    return response, nil
}

func (eh *EncounterHandlergRPC) GetEncounterByIDRpc(ctx context.Context, req *Encounters.GetEncounterByIDRequest) (*Encounters.GetEncounterByIDResponse, error) {
	log.Printf("Request received: %+v\n", req)

	encounter, err := eh.EncounterService.GetEncounterByID(req.Id)
	if err != nil {
		eh.logger.Printf("Database exception: %v", err)
		return nil, err
	}

	if encounter == nil {
		return nil, status.Errorf(codes.NotFound, "Encounter with given id not found")
	}

	return &Encounters.GetEncounterByIDResponse{
		Encounter: &Encounters.Encounter{
			Id:               encounter.ID.Hex(), // Convert ObjectID to string
			Name:             encounter.Name,
			Description:      encounter.Description,
			XpPoints:         encounter.XpPoints,
			Status:           Encounters.EncounterStatus(encounter.Status),
			Type:             Encounters.EncounterType(encounter.Type),
			Latitude:         encounter.Latitude,
			Longitude:        encounter.Longitude,
			ShouldBeApproved: encounter.ShouldBeApproved,
		},
	}, nil
}








