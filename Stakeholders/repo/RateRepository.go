package repo

import (
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