package configs

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func EnvMongoURI() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv("MONGO_URI")
}

func EnvSecretKey() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file property SECRET_JWT_KEY")
	}

	return os.Getenv("SECRET_JWT_KEY")
}
