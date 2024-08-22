package handler

import (
	TOURS "tours/proto"
	"tours/service"
)

type TourPointHandlerGRPC struct {
	TourPointService *service.TourPointService
	TOURS.UnimplementedToursServiceServer
}

func NewTourPointHandlerGRPC(tp *service.TourPointService) *TourPointHandlerGRPC {
	return &TourPointHandlerGRPC{
		TourPointService: tp,
	}
}


