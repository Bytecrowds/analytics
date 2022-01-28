package configuration

import (
	"os"
	"github.com/joho/godotenv"
)

type Config struct{
	Port string
	ConnectionString string
	Database string
}

func GetConfig() *Config {
	godotenv.Load()
	
	port := os.Getenv("PORT")
	if port == "" {
    	port = "5000"
	}

	connectionString := os.Getenv("CONNECTION_STRING")

	database := os.Getenv("DATABASE")
	if database == "" {
    	database = "testingDB"
	}

	config := Config{Port: port, ConnectionString: connectionString, Database: database}

	return &config
}



