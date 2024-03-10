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

func startServer(userHandler *handler.UserHandler, authHandler *handler.AuthHandler){
	router := mux.NewRouter()

	router.HandleFunc("/register",userHandler.RegisterUser).Methods("POST")
	router.HandleFunc("/login",authHandler.Login).Methods("POST")
	router.HandleFunc("/userProfile/{id}",userHandler.GetProfile).Methods("GET")
	router.HandleFunc("/updateProfile",userHandler.UpdateProfile).Methods("PUT")
	
	println("Server starting")
	log.Fatal(http.ListenAndServe(":8080",router))
}

func main() {

	database := initDB()
	if database == nil {
		print("FAILED TO CONNECT TO DB")
		return
	}
	userRepo:=&repo.UserRepository{DatabaseConnection: database}
	userService:=&service.UserService{UserRepo: userRepo}
	userHandler := &handler.UserHandler{UserService: userService}
	authRepo := &repo.AuthRepository{DatabaseConnection: database}
	authService :=&service.AuthService{AuthRepo: authRepo}
	authHandler := &handler.AuthHandler{AuthService: authService}
	startServer(userHandler,authHandler)
}