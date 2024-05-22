package main

import (
	"context"
	"example/gateway/config"
	"example/gateway/proto/greeter"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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
	
	conn2 , err := grpc.DialContext(
		context.Background(),
		cfg.EncountersServiceAddress,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	conn3 , err := grpc.DialContext(
		context.Background(),
		cfg.StakeholdersServiceAddress,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	
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

	gwServer := &http.Server{
		Addr:    cfg.Address,
		Handler: gwmux,
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
