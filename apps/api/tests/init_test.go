package tests

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"logger"
	"onlineshop/api/internal/handler"
	mock_handler "onlineshop/api/internal/handler/mocks"
	"onlineshop/database"
	"repositories"
	"services"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

var mongoURI, dbName string

func init() {
	mongoURI = os.Getenv("MONGO_URI")
	dbName = os.Getenv("DB_NAME")
}

type AppTestSuite struct {
	suite.Suite

	db           *database.Mongo
	api          *handler.Handler
	repositories *repositories.Repositories
	services     *services.Services
	app          *fiber.App
}

func TestAPISuite(t *testing.T) {
	if testing.Short() {
		t.Skipf("skip e2e test")
	}

	suite.Run(t, new(AppTestSuite))
}

func (s *AppTestSuite) SetupSuite() {
	s.setupDeps()
}

func (s *AppTestSuite) TearDownSuite() {
	s.db.Close(context.Background()) //nolint:errcheck
}

func (s *AppTestSuite) TearDownSubTest(suiteName, testName string) {
	logger.Get().Sugar().Infof("running: %s", testName)
}

func (s *AppTestSuite) setupDeps() {
	logger.NewLogger(logger.Config{
		EnableStacktrace: true,
	})
	logger.Get().Info("Booting e2e test")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	mongo, err := database.Connect(ctx, mongoURI, dbName)
	if err != nil {
		s.FailNow("failed to connect to mongodb", err)
		return
	}

	repos := repositories.NewRepositories(mongo, nil, nil)

	svc := services.NewServices(repos, mongo)
	mockUploader := mock_handler.NewMockImageUploader(gomock.NewController(s.T()))
	apiHandler := handler.NewHandler(repos, svc, mockUploader)

	app := fiber.New(fiber.Config{
		Immutable:    true,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			logger.Get().Error("error in e2e test", zap.Error(err))
			return c.SendStatus(http.StatusInternalServerError)
		},
	})

	apiHandler.RegisterRoutes(app, "abcd")

	s.app = app
	s.db = mongo
	s.services = &svc
	s.api = apiHandler
	s.repositories = &repos
}
