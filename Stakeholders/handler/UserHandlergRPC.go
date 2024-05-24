package handler

import (
	"context"
	"fmt"
	"log"
	"stakeholders/model"
	stakeholders "stakeholders/proto"
	"stakeholders/service"

	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandlergRPC struct {
	UserService *service.UserService
	AuthService *service.AuthService
	stakeholders.UnimplementedStakeholderServiceServer
}

func NewUserHandlergRPC(us  *service.UserService, as *service.AuthService) *UserHandlergRPC {
	return &UserHandlergRPC{
		UserService: us,
		AuthService: as,
	}
}

func (uh *UserHandlergRPC) RegistrationRpc(ctx context.Context , req *stakeholders.RegistrationRequest) (*stakeholders.RegistrationResponse,error) {

    tracer := otel.Tracer("stakeholder-service")
	_, span := tracer.Start(ctx, "Register-User")
	defer span.End()

	registration := model.Registration {
		Username: req.Username,
        Password: req.Password,
        Email: req.Email,
        Name: req.Name,
        Surname: req.Surname,
        Role: req.Role,
	}

	token := uh.AuthService.GenerateUniqueVerificationToken()
    item := false

	err := uh.UserService.Registration(&registration, &token, &item)
    if err != nil {
        return nil, fmt.Errorf("error while registering a new user: %v", err)
    }

	err = uh.AuthService.SendVerificationMail(&registration, token)
	if err != nil {
		return nil, fmt.Errorf("Error whie sending a email: %v", err)
	}

	return &stakeholders.RegistrationResponse{
        Message: "User registered successfully",
    }, nil

}

func (uh *UserHandlergRPC) GetProfileRpc(ctx context.Context , req *stakeholders.GetProfileRequest) (*stakeholders.GetProfileResponse,error) {
	tracer := otel.Tracer("stakeholder-service")
	_, span := tracer.Start(ctx, "Get-Profile")
	defer span.End()
	userId := req.Id;
	person ,err := uh.UserService.GetPersonByUserId(&userId)

	if err != nil {
        return nil, status.Errorf(codes.NotFound, "User not found: %v", err)
    }

	return &stakeholders.GetProfileResponse{
        Person: &stakeholders.Person{
			Id: uint64(person.ID),
			UserId: uint32(person.UserID),
			Name: person.Name,
			Surname: person.Surname,
			ProfileImage: person.Image,
			Email: person.Email,
			Bio: person.Bio,
			Quote: person.Quote,
        },
    },nil
}

func (uh *UserHandlergRPC) UpdateProfileRpc(ctx context.Context, req *stakeholders.UpdateProfileRequest) (*stakeholders.UpdateProfileResponse, error) {
	tracer := otel.Tracer("stakeholder-service")
	_, span := tracer.Start(ctx, "Update profile")
	defer span.End()
	log.Printf(req.String())
	//log.Printf("Received person: %+v", person)
    // Convert gRPC Person message to the model.Person struct
    modelPerson := &model.Person{
		
        UserID:       uint(req.GetUserId()),
        Name:         req.GetName(),
        Surname:      req.GetSurname(),
        Image:        req.GetProfileImage(),
        Email:        req.GetEmail(),
        Bio:          req.GetBio(),
        Quote:        req.GetQuote(),
    }
    
    // Call the service to update the profile
    updatedPerson, err := uh.UserService.UpdateProfile(modelPerson)
    if err != nil {
        return nil, status.Errorf(codes.NotFound, "Error updating profile: %v", err)
    }
    
    // Convert the updated model.Person back to gRPC Person message
    responsePerson := &stakeholders.Person{
        Id:           uint64(updatedPerson.ID),
        UserId:       uint32(updatedPerson.UserID),
        Name:         updatedPerson.Name,
        Surname:      updatedPerson.Surname,
        ProfileImage: updatedPerson.Image,
        Email:        updatedPerson.Email,
        Bio:          updatedPerson.Bio,
        Quote:        updatedPerson.Quote,
    }

    return &stakeholders.UpdateProfileResponse{Person: responsePerson}, nil
}

