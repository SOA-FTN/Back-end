package main

import (
	"log"
	"net/http"
	"stakeholders/handler"
	"stakeholders/model"
	"stakeholders/repo"
	"stakeholders/service"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {
	connection_url := "user=postgres password=super dbname=SOA port=5432 sslmode=disable"
	database,err := gorm.Open(postgres.Open(connection_url),&gorm.Config{})

	if err != nil {
		print(err)
		return nil
	}
	database.AutoMigrate(&model.User{})
	database.AutoMigrate(&model.Person{})
	return database
}

func startServer(handler *handler.UserHandler){
	router := mux.NewRouter()

	router.HandleFunc("/register",handler.RegisterUser).Methods("POST")
	router.HandleFunc("/login",handler.Login).Methods("POST")
	
	println("Server starting")
	log.Fatal(http.ListenAndServe(":8080",router))
}

func main() {

	database := initDB()
	if database == nil {
		print("FAILED TO CONNECT TO DB")
		return
	}
	repo:=&repo.UserRepository{DatabaseConnection: database}
	service:=&service.UserService{UserRepo: repo}
	handler := &handler.UserHandler{UserService: service}

	startServer(handler)
}