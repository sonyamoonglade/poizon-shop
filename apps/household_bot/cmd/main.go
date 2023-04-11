package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"household_bot/config"
	"household_bot/internal/catalog"
	"household_bot/internal/telegram/bot"
	"household_bot/internal/telegram/handler"
	"household_bot/internal/telegram/router"
	"logger"
	"onlineshop/database"
	"repositories"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	configPath, logsPath, production, strict := readCmdArgs()

	if err := logger.NewLogger(logger.Config{
		Out:              []string{logsPath},
		Strict:           strict,
		Production:       production,
		EnableStacktrace: false,
	}); err != nil {
		return fmt.Errorf("error instantiating logger: %w", err)
	}

	cfg, err := config.ReadConfig(configPath)
	if err != nil {
		return fmt.Errorf("can't read config: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	mongo, err := database.Connect(ctx, cfg.Database.URI, cfg.Database.Name)
	if err != nil {
		return fmt.Errorf("error connecting to mongo: %w", err)
	}

	repos := repositories.NewRepositories(mongo, nil, nil)

	tgBot, err := bot.NewBot(bot.Config{
		Token: cfg.Bot.Token,
	})
	if err != nil {
		return fmt.Errorf("error creating telegram bot: %w", err)
	}
	catalogProvider := catalog.NewProvider()
	tgHandler := handler.NewHandler(tgBot, repos.Rate, repos.HouseholdCategory, catalogProvider)
	tgRouter := router.NewRouter(tgBot.GetUpdates(),
		tgHandler,
		nil,
		cfg.Bot.HandlerTimeout)

	if err := tgRouter.Bootstrap(); err != nil {
		return err
	}

	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt, syscall.SIGINT)

	// Graceful shutdown
	<-exitChan

	return mongo.Close(context.Background())
}

func readCmdArgs() (string, string, bool, bool) {
	production := flag.Bool("production", false, "if logger should write to file")
	logsPath := flag.String("logs-path", "", "where log file is")
	strict := flag.Bool("strict", false, "if logger should log only warn+ logs")
	configPath := flag.String("config-path", "", "where config file is")

	flag.Parse()

	// Critical for app if not specified
	if *configPath == "" {
		panic("config path is not provided")
	}

	// Naked return, see return variable names
	return *configPath, *logsPath, *production, *strict
}
