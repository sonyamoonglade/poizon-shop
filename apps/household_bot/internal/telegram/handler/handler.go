package handler

import (
	"context"
	"errors"

	"domain"
	"household_bot/internal/catalog"
	"repositories"
	"services"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	ErrInvalidState      = errors.New("invalid state")
	ErrInvalidPriceInput = errors.New("invalid price input")
)

type RateProvider interface {
	GetYuanRate(ctx context.Context) (float64, error)
}
type Bot interface {
	Send(c tg.Chattable) (tg.Message, error)
	CleanRequest(c tg.Chattable) error
	SendMediaGroup(c tg.MediaGroupConfig) ([]tg.Message, error)
}

type handler struct {
	bot             Bot
	rateProvider    RateProvider
	categoryRepo    repositories.HouseholdCategory
	customerRepo    repositories.HouseholdCustomer
	orderService    services.Order[domain.HouseholdOrder]
	catalogProvider *catalog.Provider
}

func NewHandler(b Bot,
	rp RateProvider,
	repos repositories.Repositories,
	catalogProvider *catalog.Provider,
	orderService services.Order[domain.HouseholdOrder],
) *handler {
	return &handler{
		bot:             b,
		rateProvider:    rp,
		categoryRepo:    repos.HouseholdCategory,
		customerRepo:    repos.HouseholdCustomer,
		catalogProvider: catalogProvider,
		orderService:    orderService,
	}
}

func (h *handler) AnswerCallback(c *tg.CallbackQuery) error {
	return h.cleanSend(tg.NewCallback(c.ID, ""))
}
