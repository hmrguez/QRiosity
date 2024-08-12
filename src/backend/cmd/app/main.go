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

	// Access the environment variables
	env := os.Getenv("ENV")
	llm := os.Getenv("LLM")

	log.Printf("Running in %s environment with LLM: %s", env, llm)

	localConfig := config.NewLocalConfig()
	server.StartServer(*localConfig)
}
