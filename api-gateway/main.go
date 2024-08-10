package main

import (
	"bytes"
	"context"
	"example/gateway/config"
	"example/gateway/proto/greeter"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang-jwt/jwt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func allowCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func validateToken(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		allowedEndpoints := map[string]bool{
            "/api/auth/login":                   true,
            "/api/stakeholders/registration":    true,
			"/api/encounter/createEncounter":    true,
			"/api/encounter/getAll":    true,
			"/api/encounter/deleteEncounter": true,
			"/api/stakeholders/getProfile":true
        }

		if allowed, ok := allowedEndpoints[r.URL.Path]; ok && allowed {
            next.ServeHTTP(w, r)
            return
        }

        tokenString := r.Header.Get("Authorization")
        if tokenString == "" {
            http.Error(w, "Missing token", http.StatusUnauthorized)
            return
        }

        tokenString = tokenString[len("Bearer "):]

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
            }
            // Replace this with your secret key
            return []byte("explorer_secret_key"), nil
        })

        if err != nil || !token.Valid {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        next.ServeHTTP(w, r)
    })
}

func ForwardToAWSAPI(method string, awsAPIURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}


		req, err := http.NewRequest(method, awsAPIURL, bytes.NewBuffer(body))
		if err != nil {
			http.Error(w, "Failed to create request", http.StatusInternalServerError)
			return
		}

		req.Header.Set("Content-Type", r.Header.Get("Content-Type"))
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "Failed to forward request", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Failed to read response body", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(resp.StatusCode)
		w.Write(responseBody)
	}
}

func main() {
	cfg := config.GetConfig()
	
	conn, err := grpc.DialContext(
		context.Background(),
		cfg.GreeterServiceAddress,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}
	/*
	conn2 , err := grpc.DialContext(
		context.Background(),
		cfg.EncountersServiceAddress,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	*/
	/*
	conn3 , err := grpc.DialContext(
		context.Background(),
		cfg.StakeholdersServiceAddress,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	*/
	if err != nil {
		log.Fatalln("Failed to dial Encounter server:", err)
	}

	gwmux := runtime.NewServeMux()
	
	client := greeter.NewGreeterServiceClient(conn)
	err = greeter.RegisterGreeterServiceHandlerClient(
		context.Background(), 
		gwmux,
		client,
	)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}
	/*
	clientEncounter := greeter.NewEncounterServiceClient(conn2)
	err = greeter.RegisterEncounterServiceHandlerClient(
		context.Background(),
		gwmux,
		clientEncounter,
	)
	if err != nil {
		log.Fatalln("Failed to register EncounterService gateway:", err)
	}
	
	clientExecution := greeter.NewEncounterExecutionServiceClient(conn2)
	err = greeter.RegisterEncounterExecutionServiceHandlerClient(
		context.Background(),
		gwmux,
		clientExecution,
	)
	*/
	/*
	clientStakeholder := greeter.NewStakeholderServiceClient(conn3)
	err = greeter.RegisterStakeholderServiceHandlerClient(
		context.Background(),
		gwmux,
		clientStakeholder,
	)

	clientAuthentication := greeter.NewAuthServiceClient(conn3)
	err = greeter.RegisterAuthServiceHandlerClient(
		context.Background(),
		gwmux,
		clientAuthentication,
	)
	*/
	mux := http.NewServeMux()
	mux.Handle("/api/", gwmux)
	mux.HandleFunc("/api/encounter/createEncounter", ForwardToAWSAPI("POST", "https://jfhz3ftx19.execute-api.us-east-1.amazonaws.com/test/encounters"))
	mux.HandleFunc("/api/encounter/getAll", ForwardToAWSAPI("GET", "https://jfhz3ftx19.execute-api.us-east-1.amazonaws.com/test/encounters"))
	mux.HandleFunc("/api/encounter/deleteEncounter", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("ID")
		if id == "" {
			http.Error(w, "Missing ID query parameter", http.StatusBadRequest)
			return
		}
		// Append the ID query parameter to the AWS API URL
		ForwardToAWSAPI("DELETE", fmt.Sprintf("https://jfhz3ftx19.execute-api.us-east-1.amazonaws.com/newproduction/singleEncounter?ID=%s", id))(w, r)
	})
	////////////////////////////////////
	mux.HandleFunc("/api/stakeholders/registration",ForwardToAWSAPI("POST","https://jfhz3ftx19.execute-api.us-east-1.amazonaws.com/test/stakeholders"))
	mux.HandleFunc("/api/auth/login",ForwardToAWSAPI("POST","https://jfhz3ftx19.execute-api.us-east-1.amazonaws.com/test/stakeholders/login"))
	mux.HandleFunc("/api/stakeholders/getProfile",func(w http.ResponseWriter, r *http.Request){
		id := r.URL.Query().Get("ID")
		if id == "" {
			// Ako ID nije prisutan, vratimo grešku
			http.Error(w, "Missing ID query parameter", http.StatusBadRequest)
			return
		}
		ForwardToAWSAPI("GET", fmt.Sprintf("https://jfhz3ftx19.execute-api.us-east-1.amazonaws.com/test/stakeholders?ID=%s", id))(w, r)
	})
	gwServer := &http.Server{
		Addr:    cfg.Address,
		Handler: allowCORS(validateToken(mux)),
	}

	go func() {
		if err := gwServer.ListenAndServe(); err != nil {
			log.Fatal("server error: ", err)
		}
	}()

	stopCh := make(chan os.Signal)
	signal.Notify(stopCh, syscall.SIGTERM)

	<-stopCh

	if err = gwServer.Close(); err != nil {
		log.Fatalln("error while stopping server: ", err)
	}
}
