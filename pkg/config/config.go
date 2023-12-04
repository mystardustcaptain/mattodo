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
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}
