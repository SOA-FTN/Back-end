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