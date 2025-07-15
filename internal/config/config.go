package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port  string
	Debug bool
	DB    DatabaseConfig
}

type DatabaseConfig struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     int
	SSLMode  string
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

	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))

	return Config{
		Port:  port,
		Debug: debug,
		DB: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "postgres"),
			Port:     dbPort,
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
