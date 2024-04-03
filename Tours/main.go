package main

import (
	"log"
	"net/http"
	"os"
	"tours/handler"
	"tours/model"
	"tours/repo"
	"tours/service"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {
	database, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})

	if err != nil {
		print(err)
		return nil
	}

	database.AutoMigrate(&model.Tour{}, &model.TourPoint{}, &model.TourReview{})
	return database
}

func main() {
	database := initDB()
	if database == nil {
		print("FAILED TO CONNECT TO DB")
		return
	}

	// Initialize repositories
	tourRepo := repo.NewTourRepository(database)
	tourPointRepo := repo.NewTourPointRepository(database)
	tourReviewRepo := repo.NewTourReviewRepository(database)

	// Initialize services
	tourService := service.NewTourService(tourRepo, tourPointRepo)
	tourPointService := service.NewTourPointService(tourPointRepo)
	tourReviewService := service.NewTourReviewService(tourReviewRepo)

	// Initialize handlers
	tourHandler := handler.NewTourHandler(tourService)
	tourPointHandler := handler.NewTourPointHandler(tourPointService)
	tourReviewHandler := handler.NewTourReviewHandler(tourReviewService)
	// Set up routes
	router := mux.NewRouter()
	router.HandleFunc("/createTour", tourHandler.CreateTourHandler).Methods("POST")
	router.HandleFunc("/createTourPoint", tourPointHandler.CreateTourPointHandler).Methods("POST")
	router.HandleFunc("/toursByUserId", tourHandler.GetToursByUserIDHandler).Methods("GET")
	router.HandleFunc("/tourPoints", tourPointHandler.GetTourPointsByTourIDHandler).Methods("GET")
	router.HandleFunc("/updateTour", tourHandler.UpdateTourHandler).Methods("PUT")
	router.HandleFunc("/publishTour/{tourID}", tourHandler.PublishTourHandler).Methods("PUT")
	router.HandleFunc("/archiveTour/{tourID}", tourHandler.ArchiveTourHandler).Methods("PUT")
	router.HandleFunc("/createTourReview", tourReviewHandler.CreateTourReviewHandler).Methods("POST")
	router.HandleFunc("/getPublishedTours", tourHandler.GetPublishedToursHandler).Methods("GET")
	router.HandleFunc("/tourReview/getTourReviewsByTourID", tourReviewHandler.GetTourReviewsByTourIDHandler).Methods("GET")

	// Start the server
	log.Println("Server started on port 8081")
	log.Fatal(http.ListenAndServe(":8081", router))
}
