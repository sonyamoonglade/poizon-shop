package telegram

import (
	"context"
	"errors"

	"clothes_bot/internal/catalog"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"repositories"
	"services"
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
	b               Bot
	customerRepo    repositories.Customer
	rateProvider    RateProvider
	orderService    services.Order
	catalogProvider *catalog.CatalogProvider
}

func NewHandler(bot Bot,
	customerRepo repositories.Customer,
	orderService services.Order,
	rateProvider RateProvider,
	catalogProvider *catalog.CatalogProvider) *handler {
	return &handler{
		b:               bot,
		customerRepo:    customerRepo,
		orderService:    orderService,
		catalogProvider: catalogProvider,
		rateProvider:    rateProvider,
	}
}

func (h *handler) AnswerCallback(callbackID string) error {
	return h.cleanSend(tg.NewCallback(callbackID, ""))
}

func (h *handler) HandleError(ctx context.Context, err error, m tg.Update) {
	if errors.Is(err, ErrInvalidPriceInput) {
		h.sendMessage(m.FromChat().ID, "Неправильный формат ввода")
		return
	}
	h.b.Send(tg.NewMessage(m.FromChat().ID, "Извини, я не понимаю тебя :("))
}
