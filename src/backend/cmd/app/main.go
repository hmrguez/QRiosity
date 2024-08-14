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

	//result, err := authService.Login("zealot.algo@gmail.com", "Test@1234")
	//if err != nil {
	//	fmt.Println("Login failed:", err)
	//	return
	//}

	//result, err := authService.SignUp("zealotZ34L0T@gmail.com", "Nick", "Test@1234")
	//if err != nil {
	//	fmt.Println("Sign up failed:", err)
	//	return
	//}

	//result, err := authService.ConfirmSignUp("zealot.algo@gmail.com", "842968")
	//if err != nil {
	//	fmt.Println("Congirmation failed", err)
	//	return
	//}

	//result, err := authService.ResendConfirmationCode("zealotZ34L0T@gmail.com")
	//if err != nil {
	//	fmt.Println("Resend confirmation code failed:", err)
	//	return
	//}
	//
	//fmt.Println("Sign up successful:", result)
}
