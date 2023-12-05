package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	// Load environment variables from .env file
	// if running locally
	if os.Getenv("DOCKER_ENV_SET") != "true" {
		log.Printf("Non-docker environment detected, loading .env file")

		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		log.Printf("Environment variables loaded from .env file")
	} else {
		log.Printf("Docker environment detected, skipping .env file load")
	}
}
