package main

import (
	"context"
	"encounters/handler"
	"encounters/repo"
	"encounters/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

/*
func initDB() *gorm.DB {
	//connection_url := "user=postgres password=super dbname=SOA-encounters port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})

	if err != nil {
		print(err)
		return nil
	}
	database.AutoMigrate(&model.Encounter{}, &model.EncounterExecution{})
	return database
}
*/
func main() {
	//database := initDB()
	//if database == nil {
		//print("FAILED TO CONNECT TO DB")
		//return
	//}
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8083"
	}

	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	logger := log.New(os.Stdout, "[product-api] ", log.LstdFlags)
	storeLogger := log.New(os.Stdout, "[encounter-store] ", log.LstdFlags)
	executionLogger := log.New(os.Stdout, "[execution-store]", log.LstdFlags)
	//ENCOUNTERS
	encounterStore, err := repo.NewEncounterRepository(timeoutContext, storeLogger)
	if err != nil {
		logger.Fatal(err)
	}
	defer encounterStore.Disconnect(timeoutContext)
	encounterStore.Ping()
	//ENCOUNTER EXECUTION
	encounterExecutionStore ,err := repo.NewEncounterExecutionRepository(timeoutContext,executionLogger)
	if err != nil {
		logger.Fatal(err)
	}
	defer encounterExecutionStore.Disconnect(timeoutContext)
	encounterExecutionStore.Ping()

	encounterService := service.NewEncounterService(logger,encounterStore)
	encounterHandler := handler.NewEncounterHandler(logger,encounterService)

	encounteExecutionService := service.NewEncounterExecutionService(logger,encounterExecutionStore)
	encounterExecutionHandler := handler.NewEncounterExecutionHandler(logger,encounteExecutionService)

	router := mux.NewRouter()
	router.Use(encounterHandler.MiddlewareContentTypeSet)

	createEncounterRouter :=router.Methods(http.MethodPost).Subrouter()
	createEncounterRouter.HandleFunc("/createEncounter",encounterHandler.CreateEncounterHandler)
	createEncounterRouter.Use(encounterHandler.MIddlewareEncounterDeserialization)

	getAllEncountersRouter := router.Methods(http.MethodGet).Subrouter();
	getAllEncountersRouter.HandleFunc("/getEncounters", encounterHandler.GetAllEncountersHandler)

	getEncounterRouter := router.Methods(http.MethodGet).Subrouter();
	getEncounterRouter.HandleFunc("/getEncounter/{encounterId}", encounterHandler.GetEncounterByIDHandler)

	createExecutionRouter :=router.Methods(http.MethodPost).Subrouter()
	createExecutionRouter.HandleFunc("/createEncounterExecution", encounterExecutionHandler.CreateEncounterExecutionHandler)
	createExecutionRouter.Use(encounterExecutionHandler.MIddlewareEncounterExecutionDeserialization)
	
	getAllExecutionsRouter :=router.Methods(http.MethodGet).Subrouter();
	getAllExecutionsRouter.HandleFunc("/getEncounterExecutions", encounterExecutionHandler.GetAllEncounterExecutionsHandler)

	getActiveEncounterRouter := router.Methods(http.MethodGet).Subrouter()
	getActiveEncounterRouter.HandleFunc("/activeEncounterByUserId/{userId}", encounterExecutionHandler.GetEncounterExecutionByUserIDAndNotCompletedHandler)

	completeExecutionRouter := router.Methods(http.MethodGet).Subrouter()
	completeExecutionRouter.HandleFunc("/completeExecution/{userId}", encounterExecutionHandler.UpdateEncounterExecutionHandler)

	cors := gorillaHandlers.CORS(gorillaHandlers.AllowedOrigins([]string{"*"}))

	server := http.Server{
		Addr:         ":" + port,
		Handler:      cors(router),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	logger.Println("Server listening on port", port)
	//Distribute all the connections to goroutines
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()


	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)
	signal.Notify(sigCh, os.Kill)

	sig := <-sigCh
	logger.Println("Received terminate, graceful shutdown", sig)

	//Try to shutdown gracefully
	if server.Shutdown(timeoutContext) != nil {
		logger.Fatal("Cannot gracefully shutdown...")
	}
	logger.Println("Server stopped")

	

	/*
	// Initialize repositories
	encounterRepo := repo.NewEncounterRepository(database)
	encounterExecutionRepo := repo.NewEncounterExecutionRepository(database)

	// Initialize services
	encounterService := service.NewEncounterService(encounterRepo)
	encounterExecutionService := service.NewEncounterExecutionService(encounterExecutionRepo)

	// Initialize handlers
	encounterHandler := handler.NewEncounterHandler(encounterService)
	encounterExecutionHandler := handler.NewEncounterExecutionHandler(encounterExecutionService)

	// Set up routes
	router := mux.NewRouter()
	router.HandleFunc("/createEncounter", encounterHandler.CreateEncounterHandler).Methods("POST")
	router.HandleFunc("/getEncounters", encounterHandler.GetAllEncountersHandler).Methods("GET")
	router.HandleFunc("/getEncounterExecutions", encounterExecutionHandler.GetAllEncounterExecutionsHandler).Methods("GET")
	router.HandleFunc("/createEncounterExecution", encounterExecutionHandler.CreateEncounterExecutionHandler).Methods("POST")
	router.HandleFunc("/activeEncounterByUserId/{userId}", encounterExecutionHandler.GetEncounterExecutionByUserIDAndNotCompletedHandler).Methods("GET")
	router.HandleFunc("/getEncounter/{encounterId}", encounterHandler.GetEncounterByIDHandler).Methods("GET")
	router.HandleFunc("/completeExecution/{userId}", encounterExecutionHandler.UpdateEncounterExecutionHandler).Methods("GET")

	// Start the server
	log.Println("Server started on port 8083")
	log.Fatal(http.ListenAndServe(":8083", router))
	*/
}


