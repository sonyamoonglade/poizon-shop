package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

var ErrConfigNoExist = errors.New("config file doesn't exist")

type AppConfig struct {
	Database struct {
		// Connection string
		URI string
		// Name of database
		Name string
	}

	Bot struct {
		// Tg bot token
		Token string
		// HandlerTimeout
		HandlerTimeout time.Duration
	}

	Redis struct {
		Addr string
	}
}

func ReadConfig(path string) (AppConfig, error) {
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return AppConfig{}, ErrConfigNoExist
		}
		return AppConfig{}, err
	}
	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if err != nil {
		return AppConfig{}, err
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		return AppConfig{}, fmt.Errorf("missing MONGO_URI env")
	}

	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		return AppConfig{}, fmt.Errorf("missing DB_NAME env")
	}

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		return AppConfig{}, fmt.Errorf("missing BOT_TOKEN env")
	}

	handlerTimeout := viper.GetInt("telegram.handler_timeout")
	if handlerTimeout == 0 {
		return AppConfig{}, fmt.Errorf("missing telegram.handler_timeout")
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		return AppConfig{}, fmt.Errorf("missing REDIS_ADDR env")
	}

	return AppConfig{
		Database: struct {
			URI  string
			Name string
		}{
			URI:  mongoURI,
			Name: dbname,
		},
		Bot: struct {
			Token          string
			HandlerTimeout time.Duration
		}{
			Token:          botToken,
			HandlerTimeout: time.Duration(handlerTimeout) * time.Second,
		},
		Redis: struct {
			Addr string
		}{
			Addr: redisAddr,
		},
	}, nil
}
