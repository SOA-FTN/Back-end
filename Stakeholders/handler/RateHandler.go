package handler

import (
	"encoding/json"
	"net/http"
	"stakeholders/model"
	"stakeholders/service"
)

type RateHandler struct {
	RateService *service.RateService
}
//REGISTRACIJA KORISNIKA
func (rateHandler *RateHandler) RateApp(writer http.ResponseWriter, req *http.Request) {
	var rate model.Rate

	err := json.NewDecoder(req.Body).Decode(&rate)
	if err != nil {
		println("Error while parsing json")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = rateHandler.RateService.RateApp(&rate)
	if err != nil {
		println("Error while rating app")
		writer.WriteHeader(http.StatusExpectationFailed)
		return
	}

	writer.WriteHeader(http.StatusCreated)
	writer.Header().Set("Content-Type" , "application/json")
}