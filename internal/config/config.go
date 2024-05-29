package config

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		// TODO: look into this later
		log.Printf("Error loading .env file")
	}
}
