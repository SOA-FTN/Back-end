package main

import (
	"encounters/handler"
	"encounters/model"
	"encounters/repo"
	"encounters/service"
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

func main() {
	database := initDB()
	if database == nil {
		print("FAILED TO CONNECT TO DB")
		return
	}

	// Initialize repositories
	encounterRepo := repo.NewEncounterRepository(database)

	// Initialize services
	encounterService := service.NewEncounterService(encounterRepo)

	// Initialize handlers
	encounterHandler := handler.NewEncounterHandler(encounterService)

	// Set up routes
	router := mux.NewRouter()
	router.HandleFunc("/createEncounter", encounterHandler.CreateEncounterHandler).Methods("POST")

	// Start the server
	log.Println("Server started on port 8083")
	log.Fatal(http.ListenAndServe(":8083", router))
}
