package config

import "os"

type Config struct {
	Address               string
	GreeterServiceAddress string
	EncountersServiceAddress string
}

func GetConfig() Config {
	return Config{
		GreeterServiceAddress: os.Getenv("GREETER_SERVICE_ADDRESS"),
		Address:               os.Getenv("GATEWAY_ADDRESS"),
		EncountersServiceAddress: os.Getenv("ENCOUNTER_SERVICE_ADDRESS"),
	}
}
