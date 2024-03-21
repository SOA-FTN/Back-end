package service

import (
	"fmt"
	"stakeholders/model"
	"stakeholders/repo"
)

type RateService struct {
	RateRepo *repo.RateRepository
}



//Post
func(service *RateService) RateApp (rate *model.Rate) error {


	//canRate,  err := service.RateRepo.CheckIfUserCanRate(int(rate.UserId))
    //if !*canRate {
    //    return  fmt.Errorf("Korisnik je vec ocijenio aplikaciju" + err.Error())
    //}

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

//Get
func (service *RateService) GetAllRates() ([]model.Rate, error) {
    ratings, err := service.RateRepo.GetAllRates()
    if err != nil {
        return nil, fmt.Errorf("error getting all ratings: %v", err)
    }
    return ratings, nil
}