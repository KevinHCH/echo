package internal

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnv(key string) (string, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, continuing with environment variables only.")
	}

	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("environment variable %s is not set", key)
	}

	return value, nil
}
