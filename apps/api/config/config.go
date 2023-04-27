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

	Redis struct {
		Addr string
	}

	Uploader struct {
		Bucket    string
		Owner     string
		S3BaseURL string
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

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		return AppConfig{}, fmt.Errorf("missing REDIS_ADDR env")
	}

	bucket := os.Getenv("BUCKET")
	if bucket == "" {
		return AppConfig{}, fmt.Errorf("missing BUCKET env")
	}

	owner := os.Getenv("OWNER")
	if owner == "" {
		return AppConfig{}, fmt.Errorf("missing OWNER env")
	}

	s3BaseURL := os.Getenv("S3_BASE_URL")
	if s3BaseURL == "" {
		return AppConfig{}, fmt.Errorf("missing S3_BASE_URL env")
	}

	_, ok := os.LookupEnv("AWS_ACCESS_KEY_ID")
	if !ok {
		return AppConfig{}, fmt.Errorf("missing AWS_ACCESS_KEY_ID env")
	}

	_, ok = os.LookupEnv("AWS_SECRET_ACCESS_KEY")
	if !ok {
		return AppConfig{}, fmt.Errorf("missing AWS_SECRET_ACCESS_KEY env")
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
		Redis: struct {
			Addr string
		}{
			Addr: redisAddr,
		},
		Uploader: struct {
			Bucket    string
			Owner     string
			S3BaseURL string
		}{
			Bucket:    bucket,
			Owner:     owner,
			S3BaseURL: s3BaseURL,
		},
	}, nil
}
