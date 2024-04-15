package main

import (
	"context"
	"followings/handler"
	"followings/repo"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8085"
	}

	// Initialize context
	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//Initialize the logger we are going to use, with prefix and datetime for every log
	logger := log.New(os.Stdout, "[following-api] ", log.LstdFlags)
	storeLogger := log.New(os.Stdout, "[following-store] ", log.LstdFlags)

	// NoSQL: Initialize Movie Repository store
	store, err := repo.New(storeLogger)
	if err != nil {
		logger.Fatal(err)
	}
	defer store.CloseDriverConnection(timeoutContext)
	store.CheckConnection()

	//Initialize the handler and inject said logger
	followingHandler := handler.NewFollowingHandler(logger, store)

	//Initialize the router and add a middleware for all the requests
	router := mux.NewRouter()

	router.Use(followingHandler.MiddlewareContentTypeSet)

	postPersonNode := router.Methods(http.MethodPost).Subrouter()
	postPersonNode.HandleFunc("/person", followingHandler.CreatePerson)
	postPersonNode.Use(followingHandler.MiddlewarePersonDeserialization)

	followPersonNode := router.Methods(http.MethodPost).Subrouter()
	followPersonNode.HandleFunc("/followPerson", followingHandler.FollowPerson)

	getFollowRecommendationsNode := router.Methods(http.MethodGet).Subrouter()
	getFollowRecommendationsNode.HandleFunc("/getRecommendations/{username}", followingHandler.GetFollowRecommendations)

	getFollowedUser := router.Methods(http.MethodGet).Subrouter()
	getFollowedUser.HandleFunc("/getFollowedUsers/{username}", followingHandler.GetFollowedUsers)

	getUsersExcept := router.Methods(http.MethodGet).Subrouter()
	getUsersExcept.HandleFunc("/getUsersExcept/{username}", followingHandler.GetUsersExcept)

	isFollowedNode := router.Methods(http.MethodPost).Subrouter()
	isFollowedNode.HandleFunc("/isFollowing", followingHandler.IsFollowing)

	cors := gorillaHandlers.CORS(gorillaHandlers.AllowedOrigins([]string{"*"}))

	//Initialize the server
	server := http.Server{
		Addr:         ":" + port,
		Handler:      cors(router),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
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
