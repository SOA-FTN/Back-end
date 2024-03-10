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

//var jwtKey = []byte ("secret_key")
//Registracija korisnika
func (userHandler *UserHandler) RegisterUser(writer http.ResponseWriter, req *http.Request) {
	var userData struct {
		User model.User
		Person  model.Person
	}

	err := json.NewDecoder(req.Body).Decode(&userData)
	if err != nil {
		println("Error while parsing json")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = userHandler.UserService.RegisterUser(&userData.User , &userData.Person)
	if err != nil {
		println("Error while registering a new user")
		writer.WriteHeader(http.StatusExpectationFailed)
		return
	}

	writer.WriteHeader(http.StatusCreated)
	writer.Header().Set("Content-Type" , "application/json")
}
/*
//Login
func(UserHandler *UserHandler) Login(writer http.ResponseWriter, req *http.Request){
	var credentials model.Credentials
	err := json.NewDecoder(req.Body).Decode(&credentials)

	if err != nil{
		print("Error while parsing json")
		writer.WriteHeader(http.StatusBadRequest)
		return 
	}
	user,err := UserHandler.UserService.Authentication(&credentials)
	if err != nil || user.Password != credentials.Password || user.UserName != credentials.Username {
		writer.WriteHeader(http.StatusUnauthorized)
		return 
	}

	claims :=UserHandler.UserService.CreateClaims(user)
	expirationTime := time.Now().Add(time.Minute *5)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,*claims)
	tokenString,err := token.SignedString(jwtKey)
	if err != nil{
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(writer,&http.Cookie{
		Name: "token",
		Value: tokenString,
		Expires : expirationTime,
	})
	
}
*/
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