package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"stakeholders/config"
	"stakeholders/handler"
	"stakeholders/model"
	stakeholders "stakeholders/proto"
	"stakeholders/repo"
	"stakeholders/service"
	"syscall"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {
	database, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})

	if err != nil {
		print(err)
		return nil
	}

	database.AutoMigrate(&model.User{}, &model.Person{}, &model.Rate{})
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

	println("Server starting on port 8082")
	log.Fatal(http.ListenAndServe(":8082", router))
}

func main() {

	database := initDB()
	if database == nil {
		print("FAILED TO CONNECT TO DB")
		return
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
	//timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	//defer cancel()

	userRepo := &repo.UserRepository{DatabaseConnection: database}
	authRepo := &repo.AuthRepository{DatabaseConnection: database}
	userService := &service.UserService{UserRepo: userRepo}
	authService := &service.AuthService{AuthRepo: authRepo}
	userHandlergRPC := handler.NewUserHandlergRPC(userService,authService);
	authHandlergRPC := handler.NewAuthHandlergRPC(authService)
	userHandler := &handler.UserHandler{UserService: userService}
	//authRepo := &repo.AuthRepository{DatabaseConnection: database}
	//authService := &service.AuthService{AuthRepo: authRepo}
	authHandler := &handler.AuthHandler{AuthService: authService}
	rateRepo := &repo.RateRepository{DatabaseConnection: database}
	rateService := &service.RateService{RateRepo: rateRepo}
	rateHandler := &handler.RateHandler{RateService: rateService}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	stakeholders.RegisterStakeholderServiceServer(grpcServer,userHandlergRPC);
	stakeholders.RegisterAuthServiceServer(grpcServer,authHandlergRPC);

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal("server error: ", err)
		}
	}()

	stopCh := make(chan os.Signal)
	signal.Notify(stopCh, syscall.SIGTERM)

	<-stopCh

	grpcServer.Stop()



	startServer(userHandler, authHandler, rateHandler)
}
