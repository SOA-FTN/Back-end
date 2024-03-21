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
	database, err := gorm.Open(postgres.Open(connection_url), &gorm.Config{})

	if err != nil {
		print(err)
		return nil
	}
	database.AutoMigrate(&model.User{})
	database.AutoMigrate(&model.Person{})
	database.AutoMigrate(&model.Rate{})
	return database
}

func startServer(userHandler *handler.UserHandler, authHandler *handler.AuthHandler, rateHandler *handler.RateHandler) {
	router := mux.NewRouter()

	router.HandleFunc("/register", userHandler.Registration).Methods("POST", "OPTIONS")
	router.HandleFunc("/login", authHandler.Login).Methods("POST", "OPTIONS")
	router.HandleFunc("/userProfile/{id}", userHandler.GetProfile).Methods("GET", "OPTIONS")
	router.HandleFunc("/updateProfile", userHandler.UpdateProfile).Methods("PUT", "OPTIONS")
	router.HandleFunc("/rate-app", rateHandler.RateApp).Methods("POST", "OPTIONS")
	router.HandleFunc("/app-ratings", rateHandler.GetAllRates).Methods("GET", "OPTIONS")
	router.HandleFunc("/verifyEmail/{token}", userHandler.VerifyEmail).Methods("GET", "OPTIONS")
	router.HandleFunc("/all-profiles", userHandler.GetAllProfiles).Methods("GET", "OPTIONS")
	router.HandleFunc("/block-profile/{id}", userHandler.BlockUser).Methods("PUT", "OPTIONS")

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET,PUT,DELETE, OPTIONS,")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	println("Server starting")
	log.Fatal(http.ListenAndServe(":8082", router))
}

func main() {

	database := initDB()
	if database == nil {
		print("FAILED TO CONNECT TO DB")
		return
	}
	userRepo := &repo.UserRepository{DatabaseConnection: database}
	userService := &service.UserService{UserRepo: userRepo}
	userHandler := &handler.UserHandler{UserService: userService}
	authRepo := &repo.AuthRepository{DatabaseConnection: database}
	authService := &service.AuthService{AuthRepo: authRepo}
	authHandler := &handler.AuthHandler{AuthService: authService}
	rateRepo := &repo.RateRepository{DatabaseConnection: database}
	rateService := &service.RateService{RateRepo: rateRepo}
	rateHandler := &handler.RateHandler{RateService: rateService}

	startServer(userHandler, authHandler, rateHandler)
}
