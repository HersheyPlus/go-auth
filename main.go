package main

import (
	"github.com/HersheyPlus/go-auth/config"
	"log"

)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
        log.Fatal("Cannot load config:", err)
    }

}