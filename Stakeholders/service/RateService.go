package service

import (
	"stakeholders/model"
	"stakeholders/repo"
)

type RateService struct {
	RateRepo *repo.RateRepository
}



//REGISTRACIJA
func(service *RateService) RateApp (rate *model.Rate) error {

	newRate := model.Rate{
		UserId: rate.UserId,
		Rating: rate.Rating,
		Description: rate.Description,
	}
	
	err := service.RateRepo.RateApp(&newRate)
	if err != nil{
		return err
	}
	return nil
}