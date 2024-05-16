package handler

import (
	"context"
	"encounters/model"
	Encounters "encounters/proto"
	"encounters/service"
	"log"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EncounterExecutionHandlergRPC struct {
	logger                    *log.Logger
	EncounterExecutionService *service.EncounterExecutionService
	Encounters.UnimplementedEncounterExecutionServiceServer
}

func NewEncounterExecutionHandlergRPC(l *log.Logger,es *service.EncounterExecutionService) *EncounterExecutionHandlergRPC {
	return &EncounterExecutionHandlergRPC{
		logger:           l,
		EncounterExecutionService: es,
	}
}

func (eh *EncounterExecutionHandlergRPC) CreateExecutionRpc(ctx context.Context, req *Encounters.CreateExecutionRequest) (*Encounters.CreateExecutionResponse, error) {
	// Konvertujemo podatke iz zahteva u format koji očekuje servis
	encounterExecution := &model.EncounterExecution{
		UserID:         req.UserId,
		EncounterID:    req.EncounterId,
		CompletionTime: req.CompletionTime,
		IsCompleted:    req.IsCompleted,
	}

	// Pozivamo odgovarajuću funkciju servisa za kreiranje izvršenja susreta
	err := eh.EncounterExecutionService.CreateEncounterExecution(encounterExecution)
	if err != nil {
		return nil, err
	}

	// Vraćamo odgovor sa informacijama o kreiranom izvršenju susreta
	return &Encounters.CreateExecutionResponse{
		Execution: &Encounters.EncounterExecution{
			Id:             encounterExecution.ID.Hex(),
			UserId:         encounterExecution.UserID,
			EncounterId:    encounterExecution.EncounterID,
			CompletionTime: encounterExecution.CompletionTime,
			IsCompleted:    encounterExecution.IsCompleted,
		},
	}, nil
}

func(eh *EncounterExecutionHandlergRPC) GetAllExecutionsRpc (ctx context.Context, req *Encounters.GetAllExecutionsRequest) (*Encounters.GetAllExecutionsResponse,error) {
	executions, err := eh.EncounterExecutionService.GetAllEncounterExecutions()
	if err != nil {
        return nil, err
    }

	var protoExecutions [] *Encounters.EncounterExecution
	for _, execution := range executions {
		protoExecution := &Encounters.EncounterExecution{
			 Id: execution.ID.Hex(),
			 UserId: execution.UserID,
			 EncounterId: execution.EncounterID,
			 CompletionTime: execution.CompletionTime,
			 IsCompleted: execution.IsCompleted,
		}
		protoExecutions = append(protoExecutions,protoExecution);
	}
	response := &Encounters.GetAllExecutionsResponse{
		Executions: protoExecutions,
	}
	return response,nil
}

func(eh *EncounterExecutionHandlergRPC) GetExecutionByUserIDRpc (ctx context.Context, req *Encounters.GetExecutionByUserIdRequest) (*Encounters.GetExecutionByUserIdResponse,error) {
	userID, err := strconv.Atoi(req.UserId)
	execution ,err := eh.EncounterExecutionService.GetEncounterExecutionByUserIDAndNotCompleted(userID)
	if err != nil {
		eh.logger.Printf("Database exception: %v", err)
		return nil, err
	}
	if execution == nil {
		return nil, status.Errorf(codes.NotFound, "Encounter with given user id not found")
	}
	return &Encounters.GetExecutionByUserIdResponse{
		Execution: &Encounters.EncounterExecution{
			Id: execution.ID.Hex(),
			UserId: execution.UserID,
			EncounterId: execution.EncounterID,
			CompletionTime: execution.CompletionTime,
			IsCompleted: execution.IsCompleted,
		},
	},nil
}

func (eeh *EncounterExecutionHandlergRPC) UpdateExecutionRpc(ctx context.Context, req *Encounters.UpdateExecutionRequest) (*Encounters.UpdateExecutionResponse, error) {
	userID, err := strconv.Atoi(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid userID: %v", err)
	}

	// Pozivamo odgovarajuću funkciju servisa za ažuriranje izvršenja susreta
	if err := eeh.EncounterExecutionService.UpdateEncounterExecution(userID); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update encounter execution: %v", err)
	}

	// Vraćamo odgovor da je ažuriranje uspešno izvršeno
	return &Encounters.UpdateExecutionResponse{}, nil
}



