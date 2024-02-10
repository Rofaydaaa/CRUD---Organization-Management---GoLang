package util

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)

func EnvMongoURI() string {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
  
    return os.Getenv("MONGOURI")
}

func EnvPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Fatal("Port environment variable is not set, you will directed to port 8080")
	}
	return port
}

func EnvDatabaseNameMongoDB() string {
	databaseName := os.Getenv("MONGODB_DATABASE_NAME")
	if databaseName == "" {
        log.Fatal("MONGODB_DATABASE_NAME environment variable is not set")
    }
	return databaseName
}