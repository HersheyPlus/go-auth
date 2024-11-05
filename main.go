package main

import (
	"log"
	"github.com/HersheyPlus/go-auth/config"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
        log.Fatal("Cannot load config:", err)
    }
	log.Println("Config loaded successfully", config)
}