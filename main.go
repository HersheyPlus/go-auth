package main

import (
	"log"
	"github.com/HersheyPlus/go-auth/config"
	"github.com/HersheyPlus/go-auth/database"
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
    db := database.GetDB()
	log.Println("Config loaded successfully:", cfg.App.Name)
	log.Println("Database connected successfully", db.Name())
}