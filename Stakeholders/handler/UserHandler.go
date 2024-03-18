package handler

import (
	"encoding/json"
	"net/http"
	"stakeholders/model"
	"stakeholders/service"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	UserService *service.UserService
}
//REGISTRACIJA KORISNIKA
func (userHandler *UserHandler) Registration(writer http.ResponseWriter, req *http.Request) {
	var registration model.Registration

	err := json.NewDecoder(req.Body).Decode(&registration)
	if err != nil {
		println("Error while parsing json")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = userHandler.UserService.Registration(&registration)
	if err != nil {
		println("Error while registering a new user")
		writer.WriteHeader(http.StatusExpectationFailed)
		return
	}

	writer.WriteHeader(http.StatusCreated)
	writer.Header().Set("Content-Type" , "application/json")
}

//Getovanje PERSON-a

func (userHandler *UserHandler) GetProfile (writer http.ResponseWriter , req *http.Request){
	userId := mux.Vars(req)["id"]
	person ,err := userHandler.UserService.GetPersonByUserId(&userId)
	writer.Header().Set("Content-Type" , "application/json")
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(person)

}

//UPDATE PERSON
func(userHandler *UserHandler) UpdateProfile(writer http.ResponseWriter , req *http.Request){
	var person model.Person
	err := json.NewDecoder(req.Body).Decode(&person)

	if err != nil{
		print("Error while parsing json")
		writer.WriteHeader(http.StatusBadRequest)
		return 
	}
	updatedPerson,err := userHandler.UserService.UpdateProfile(&person)
	writer.Header().Set("Content-Type" , "application/json")
	if  err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(updatedPerson)
}