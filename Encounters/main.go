package main

import (
	"context"
	"encounters/config"
	"encounters/handler"
	Encounters "encounters/proto"
	"encounters/repo"
	"encounters/service"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)


func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8083"
	}

	cfg := config.GetConfig()
	//GRPC
	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		log.Fatalln(err)
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(listener)
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

	encounterService := service.NewEncounterService(logger,encounterStore)
	encounterHandler := handler.NewEncounterHandlergRPC(logger,encounterService);
	//ENCOUNTER EXECUTION
	encounterExecutionStore ,err := repo.NewEncounterExecutionRepository(timeoutContext,executionLogger)
	if err != nil {
		logger.Fatal(err)
	}
	defer encounterExecutionStore.Disconnect(timeoutContext)
	encounterExecutionStore.Ping()

	encounterExecutionService := service.NewEncounterExecutionService(logger,encounterExecutionStore)
	encounterExecutionHandler := handler.NewEncounterExecutionHandlergRPC(logger,encounterExecutionService)

	// Bootstrap gRPC server.
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	
	Encounters.RegisterEncounterServiceServer(grpcServer,encounterHandler);
	Encounters.RegisterEncounterExecutionServiceServer(grpcServer,encounterExecutionHandler)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal("server error: ", err)
		}
	}()

	stopCh := make(chan os.Signal)
	signal.Notify(stopCh, syscall.SIGTERM)

	<-stopCh

	grpcServer.Stop()
	//==============================================================
	
	router := mux.NewRouter()
	//router.Use(encounterHandler.MiddlewareContentTypeSet)

	//createEncounterRouter :=router.Methods(http.MethodPost).Subrouter()
	//createEncounterRouter.HandleFunc("/createEncounter",encounterHandler.CreateEncounterHandler)
	//createEncounterRouter.Use(encounterHandler.MIddlewareEncounterDeserialization)

	//getAllEncountersRouter := router.Methods(http.MethodGet).Subrouter();
	//getAllEncountersRouter.HandleFunc("/getEncounters", encounterHandler.GetAllEncountersHandler)

	//getEncounterRouter := router.Methods(http.MethodGet).Subrouter();
	//getEncounterRouter.HandleFunc("/getEncounter/{encounterId}", encounterHandler.GetEncounterByIDHandler)

	// createExecutionRouter :=router.Methods(http.MethodPost).Subrouter()
	// createExecutionRouter.HandleFunc("/createEncounterExecution", encounterExecutionHandler.CreateEncounterExecutionHandler)
	// createExecutionRouter.Use(encounterExecutionHandler.MIddlewareEncounterExecutionDeserialization)
	
	// getAllExecutionsRouter :=router.Methods(http.MethodGet).Subrouter();
	// getAllExecutionsRouter.HandleFunc("/getEncounterExecutions", encounterExecutionHandler.GetAllEncounterExecutionsHandler)

	// getActiveEncounterRouter := router.Methods(http.MethodGet).Subrouter()
	// getActiveEncounterRouter.HandleFunc("/activeEncounterByUserId/{userId}", encounterExecutionHandler.GetEncounterExecutionByUserIDAndNotCompletedHandler)

	// completeExecutionRouter := router.Methods(http.MethodGet).Subrouter()
	// completeExecutionRouter.HandleFunc("/completeExecution/{userId}", encounterExecutionHandler.UpdateEncounterExecutionHandler)

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
	
}


