package handler

import (
	"context"
	"stakeholders/model"
	stakeholders "stakeholders/proto"
	"stakeholders/service"

	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandlergRPC struct {
	AuthService *service.AuthService
	stakeholders.UnimplementedAuthServiceServer
}

func NewAuthHandlergRPC(as *service.AuthService) *AuthHandlergRPC {
	return &AuthHandlergRPC{
		AuthService: as,
	}
}


func (ah *AuthHandlergRPC) LoginRpc(ctx context.Context, req *stakeholders.LoginRequest) (*stakeholders.LoginResponse, error) {
    tracer := otel.Tracer("stakeholder-service")
	_, span := tracer.Start(ctx, "Login")
	defer span.End()

    credentials := &model.Credentials{
        Username: req.GetUsername(),
        Password: req.GetPassword(),
    }

    // Authenticate user
    user, err := ah.AuthService.Authentication(credentials)
    if err != nil || user.Password != credentials.Password || user.UserName != credentials.Username || !*user.IsActive {
        return nil, status.Errorf(codes.Unauthenticated, "authentication failed")
    }

    // Generate token
    token, err := ah.AuthService.GenerateToken(user)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "could not generate token: %v", err)
    }

    // Prepare response
    response := &stakeholders.LoginResponse{
        Id:          uint64(user.ID),
        AccessToken: token,
    }

    return response, nil
}

