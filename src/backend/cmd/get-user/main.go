package get_user

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

	env := os.Getenv("ENV")

	if env == "prod" {
		server.LambdaServer()
		//appConfig = config.NewProdConfig()
	} else {
		var appConfig *config.Config = config.NewLocalConfig()
		server.StartServer(*appConfig)
	}

}
