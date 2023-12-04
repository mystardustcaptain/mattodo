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
		log.Println("Non-docker environment detected, loading .env file")

		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		log.Println("Environment variables loaded from .env file")
	} else {
		log.Println("Docker environment detected, skipping .env file load")
	}
}
