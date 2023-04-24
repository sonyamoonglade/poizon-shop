package tests

import (
	"context"
	"os"
	"testing"
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"household_bot/internal/catalog"
	"household_bot/internal/telegram/handler"
	"household_bot/internal/telegram/router"
	mock_handler "household_bot/mocks"
	"logger"
	"onlineshop/database"
	"repositories"
	"services"
	"usecase"
)

var mongoURI, dbName string

func init() {
	mongoURI = os.Getenv("MONGO_URI")
	dbName = os.Getenv("DB_NAME")
}

type AppTestSuite struct {
	suite.Suite

	db               *database.Mongo
	repositories     *repositories.Repositories
	svc              *services.Services
	makeOrderUsecase *usecase.HouseholdMakeOrder

	mockBot     *mock_handler.MockBot
	tghandler   router.RouteHandler
	updatesChan <-chan tg.Update
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

	updates := make(chan tg.Update)

	ctrl := gomock.NewController(s.T())
	mockBot := mock_handler.NewMockBot(ctrl)

	svc := services.NewServices(repos, mongo)

	makeOrderUsecase := usecase.NewHouseholdMakeOrderUsecase(repos.Promocode, svc.HouseholdOrder, svc.HouseholdCustomer)

	s.Require().NoError(repos.Rate.UpdateRate(context.Background(), 11.8))

	tgHandler := handler.NewHandlerWithConstructor(handler.Constructor{
		Bot:               mockBot,
		RateProvider:      repos.Rate,
		PromocodeRepo:     repos.Promocode,
		CatalogProvider:   catalogProvider,
		OrderService:      svc.HouseholdOrder,
		CategoryService:   svc.HouseholdCategory,
		CatalogMsgService: svc.HouseholdCatalogMsg,
		CustomerService:   svc.HouseholdCustomer,
		MakeOrderUsecase:  makeOrderUsecase,
	})

	s.db = mongo
	s.updatesChan = updates
	s.tghandler = tgHandler
	s.repositories = &repos
	s.svc = &svc
	s.mockBot = mockBot
	s.makeOrderUsecase = makeOrderUsecase
}

func (s *AppTestSuite) replaceBotInHandler(bot handler.Bot) {
	s.tghandler = handler.NewHandlerWithConstructor(handler.Constructor{
		Bot:               bot,
		RateProvider:      s.repositories.Rate,
		MakeOrderUsecase:  s.makeOrderUsecase,
		PromocodeRepo:     s.repositories.Promocode,
		CatalogProvider:   catalog.NewProvider(),
		OrderService:      s.svc.HouseholdOrder,
		CategoryService:   s.svc.HouseholdCategory,
		CatalogMsgService: s.svc.HouseholdCatalogMsg,
		CustomerService:   s.svc.HouseholdCustomer,
	})
}
