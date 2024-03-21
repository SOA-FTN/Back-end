package handler

import (
	"encoding/json"
	"net/http"
	"stakeholders/model"
	"stakeholders/service"
)

type AuthHandler struct {
	AuthService *service.AuthService
}

//Login
func(AuthHandler *AuthHandler) Login(writer http.ResponseWriter, req *http.Request){
	var credentials model.Credentials
	err := json.NewDecoder(req.Body).Decode(&credentials)

	if err != nil{
		print("Error while parsing json")
		writer.WriteHeader(http.StatusBadRequest)
		return 
	}
	user,err := AuthHandler.AuthService.Authentication(&credentials)
	if err != nil || user.Password != credentials.Password || user.UserName != credentials.Username || !*user.IsActive {
		writer.WriteHeader(http.StatusUnauthorized)
		return 
	}

	token, err := AuthHandler.AuthService.GenerateToken(user)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := map[string]interface{} {
		"id" : user.ID,
		"accessToken":token,
	}

    writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(writer).Encode(response)
}