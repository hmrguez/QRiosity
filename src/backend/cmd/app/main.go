package main

import (
	"backend/internal/config"
	"backend/internal/server"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file", err)
	}

	var appConfig *config.Config

	env := os.Getenv("ENV")

	if env == "prod" {
		appConfig = config.NewProdConfig()
	} else {
		appConfig = config.NewLocalConfig()
	}

	server.StartServer(*appConfig)
}
