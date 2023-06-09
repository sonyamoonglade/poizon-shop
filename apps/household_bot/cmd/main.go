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

	"domain"
	"go.uber.org/zap"
	"household_bot/config"
	"household_bot/internal/catalog"
	"household_bot/internal/telegram/bot"
	"household_bot/internal/telegram/handler"
	"household_bot/internal/telegram/router"
	"logger"
	"onlineshop/database"
	"redis"
	"repositories"
	"services"
	"usecase"
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
	initialCatalog, err := repos.HouseholdCategory.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("get initial catalog: %w", err)
	}
	catalogProvider.Load(initialCatalog)

	svc := services.NewServices(repos, mongo)
	makeOrderUsecase := usecase.NewHouseholdMakeOrderUsecase(repos.Promocode, svc.HouseholdOrder, svc.HouseholdCustomer)
	tgHandler := handler.NewHandlerWithConstructor(handler.Constructor{
		Bot:               tgBot,
		RateProvider:      repos.Rate,
		PromocodeRepo:     repos.Promocode,
		CatalogProvider:   catalogProvider,
		OrderService:      svc.HouseholdOrder,
		CategoryService:   svc.HouseholdCategory,
		CatalogMsgService: svc.HouseholdCatalogMsg,
		CustomerService:   svc.HouseholdCustomer,
		MakeOrderUsecase:  makeOrderUsecase,
	},
	)

	tgRouter := router.NewRouter(
		tgBot.GetUpdates(),
		tgHandler,
		repos.HouseholdCustomer,
		cfg.Bot.HandlerTimeout,
	)

	client := redis.NewClient(cfg.Redis.Addr)
	bus := redis.NewBus[[]domain.HouseholdCategory](client)

	redCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	onCatalogUpdate := func(items []domain.HouseholdCategory) error {
		catalogProvider.Load(items)
		return nil
	}
	redisErrorHandler := func(topic string, err error) {
		logger.Get().Error("redis error", zap.Error(err))
	}
	go bus.SubscribeToTopicWithCallback(
		redCtx,
		redis.HouseholdCatalogTopic,
		onCatalogUpdate,
		redisErrorHandler,
	)

	go bus.SubscribeToTopicWithCallback(
		redCtx,
		redis.HouseholdWipeCatalogTopic,
		func(_ []domain.HouseholdCategory) error {
			return tgHandler.WipeCatalogs(redCtx)
		},
		redisErrorHandler,
	)

	go tgRouter.Bootstrap()

	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt, syscall.SIGINT)

	// Graceful shutdown
	<-exitChan
	logger.Get().Info("graceful shutdown")
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
