package config

import (
	"fmt"
	"os"
)

type AppConfig struct {
	Database struct {
		// Connection string
		URI string
		// Name of database
		Name string
	}
	HTTP struct {
		Port string
	}
}

func ReadConfig(path string) (AppConfig, error) {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		return AppConfig{}, fmt.Errorf("missing MONGO_URI env")
	}

	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		return AppConfig{}, fmt.Errorf("missing DB_NAME env")
	}

	port := os.Getenv("PORT")
	if dbname == "" {
		return AppConfig{}, fmt.Errorf("missing PORT env")
	}

	return AppConfig{
		Database: struct {
			URI  string
			Name string
		}{
			URI:  mongoURI,
			Name: dbname,
		},
		HTTP: struct {
			Port string
		}{
			Port: port,
		},
	}, nil
}
