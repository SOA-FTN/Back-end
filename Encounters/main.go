package main

import (
	"encounters/model"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {
	connection_url := "user=postgres password=super dbname=SOA-encounters port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(connection_url), &gorm.Config{})

	if err != nil {
		print(err)
		return nil
	}
	database.AutoMigrate(&model.Encounter{})

	return database
}

func startServer() {
	router := mux.NewRouter()

	println("Server starting")
	log.Fatal(http.ListenAndServe(":8083", router))
}

func main() {
	database := initDB()
	if database == nil {
		print("FAILED TO CONNECT TO DB")
		return
	}

	// Initialize repositories
	//tourRepo := repo.NewTourRepository(database)
	//tourPointRepo := repo.NewTourPointRepository(database)

	// Initialize services
	//tourService := service.NewTourService(tourRepo)
	//tourPointService := service.NewTourPointService(tourPointRepo)

	// Initialize handlers
	//tourHandler := handler.NewTourHandler(tourService)
	//tourPointHandler := handler.NewTourPointHandler(tourPointService)

	// Set up routes
	router := mux.NewRouter()
	//router.HandleFunc("/createTour", tourHandler.CreateTourHandler).Methods("POST")
	//router.HandleFunc("/createTourPoint", tourPointHandler.CreateTourPointHandler).Methods("POST")
	//router.HandleFunc("/toursByUserId", tourHandler.GetToursByUserIDHandler).Methods("GET")
	//router.HandleFunc("/tourPoints", tourPointHandler.GetTourPointsByTourIDHandler).Methods("GET")

	// Start the server
	log.Println("Server started on port 8083")
	log.Fatal(http.ListenAndServe(":8083", router))
}
