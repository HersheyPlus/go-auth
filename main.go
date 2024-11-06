package main

import (
	"log"

	"github.com/HersheyPlus/go-auth/config"
	"github.com/HersheyPlus/go-auth/database"
	"github.com/HersheyPlus/go-auth/server"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
        log.Fatal("Cannot load config:", err)
    }
	if err := database.ConnectDatabase(cfg); err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer database.CloseDB()
	server := server.NewServer(cfg)
	if err := server.RunServer(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
    
}