package tests

import (
	"context"
	"os"
	"testing"
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"household_bot/internal/catalog"
	"household_bot/internal/telegram/handler"
	"household_bot/internal/telegram/router"
	"logger"
	"onlineshop/database"
	"repositories"
	"services"
)

var mongoURI, dbName string

type MockBot struct {
	mock.Mock
}

func (mb *MockBot) CleanRequest(c tg.Chattable) error {
	args := mb.Called(c)
	return args.Error(0)
}

func (mb *MockBot) SendMediaGroup(c tg.MediaGroupConfig) ([]tg.Message, error) {
	args := mb.Called(c)
	return args.Get(0).([]tg.Message), args.Error(1)
}

func (mb *MockBot) Send(c tg.Chattable) (tg.Message, error) {
	args := mb.Called(c)
	return args.Get(0).(tg.Message), args.Error(1)
}

func init() {
	mongoURI = os.Getenv("MONGO_URI")
	dbName = os.Getenv("DB_NAME")
}

type AppTestSuite struct {
	suite.Suite

	db           *database.Mongo
	tgrouter     *router.Router
	tghandler    router.RouteHandler
	repositories *repositories.Repositories
	mockBot      *MockBot
	updatesChan  <-chan tg.Update
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

	catalogProvider := catalog.NewProvider()
	repos := repositories.NewRepositories(mongo, nil, nil)
	orderService := services.NewHouseholdOrderService(repos.HouseholdOrder)

	updates := make(chan tg.Update)
	mockBot := new(MockBot)
	tgHandler := handler.NewHandler(mockBot, repos.Rate, repos, catalogProvider, orderService)
	tgRouter := router.NewRouter(updates, tgHandler, repos.ClothingCustomer, time.Second*5)

	mockBot.On("Send", mock.Anything).Return(tg.Message{}, nil)

	s.Require().NoError(repos.Rate.UpdateRate(context.Background(), 11.8))

	s.db = mongo
	s.updatesChan = updates
	s.tgrouter = tgRouter
	s.tghandler = tgHandler
	s.repositories = &repos
	s.mockBot = mockBot
}
