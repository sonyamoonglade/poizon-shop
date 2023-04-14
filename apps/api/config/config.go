package config

import (
	"fmt"
	"os"
)

type AppConfig struct {
	Database struct {
		URI  string
		Name string
	}
	HTTP struct {
		Port   string
		ApiKey string
	}
}

func ReadConfig() (AppConfig, error) {
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

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return AppConfig{}, fmt.Errorf("missing API_KEY env")
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
			Port   string
			ApiKey string
		}{
			Port:   port,
			ApiKey: apiKey,
		},
	}, nil
}
