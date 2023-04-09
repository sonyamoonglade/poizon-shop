package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"
	"logger"
	"onlineshop/api/config"
	"onlineshop/api/internal/handler"
	"onlineshop/database"
	"repositories"
	"services"
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

	// HTTP api
	app := fiber.New(fiber.Config{
		Immutable: true,
		Prefork:   false,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			logger.Get().Error("error in api endpoint", zap.ByteString("url",
				ctx.Request().RequestURI()),
				zap.Error(err))
			return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
	}))

	app.Use(func(c *fiber.Ctx) error {
		logger.Get().Debug("new incoming request",
			zap.ByteString("url", c.Request().RequestURI()),
		)
		return c.Next()
	})

	hhCategoryService := services.NewHouseholdCategoryService(repos.HouseholdCategory)
	apiController := handler.NewHandler(repos, hhCategoryService)
	apiController.RegisterRoutes(app)

	if err := app.Listen(":" + cfg.HTTP.Port); err != nil {
		return err
	}
	logger.Get().Info("http api server is up")

	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt, syscall.SIGINT)
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
