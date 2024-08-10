package main

import (
	"backend/internal/config"
	"backend/internal/server"
)

func main() {
	localConfig := config.NewLocalConfig()
	server.StartServer(*localConfig)
}
