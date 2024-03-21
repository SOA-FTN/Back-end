package repo

import (
	"errors"
	"stakeholders/model"

	"gorm.io/gorm"
)

type RateRepository struct {
	DatabaseConnection *gorm.DB
}

func(rateRepo *RateRepository) RateApp(rate *model.Rate) error {
	dbResult := rateRepo.DatabaseConnection.Create(rate)
	if(dbResult.Error != nil) {
		return dbResult.Error
	}
	return nil
}

func (rateRepo *RateRepository) GetAllRates() ([]model.Rate, error) {
    var rates []model.Rate
    dbResult := rateRepo.DatabaseConnection.Select("id","user_id","rating","description").Find(&rates)
    if dbResult.Error != nil {
        return nil, dbResult.Error
    }
    
    return rates, nil
}

func (rateRepo *RateRepository) CheckIfUserCanRate(userID int) (*bool, error) {
    var rate model.Rate
    dbResult := rateRepo.DatabaseConnection.Where("user_id = ?", userID).First(&rate)
    if dbResult.Error != nil {
        if errors.Is(gorm.ErrRecordNotFound, dbResult.Error) {
            canRate := true
            return &canRate, nil
        }
        return nil, dbResult.Error
    }

    canRate := false
    return &canRate, nil
}