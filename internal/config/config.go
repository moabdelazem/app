package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port  string
	Debug bool
}

func Load() Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "4260"
	}

	debug := os.Getenv("DEBUG") == "true"

	return Config{
		Port:  port,
		Debug: debug,
	}
}
