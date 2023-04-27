package telegram

import (
	"context"
	"errors"

	"clothes_bot/internal/catalog"
	"domain"
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
	customerService services.ClothingCustomer
	rateProvider    RateProvider
	orderService    services.Order[domain.ClothingOrder]
	promocodeRepo   repositories.Promocode
	catalogProvider *catalog.CatalogProvider
}

func NewHandler(bot Bot,
	customerService services.ClothingCustomer,
	orderService services.Order[domain.ClothingOrder],
	rateProvider RateProvider,
	catalogProvider *catalog.CatalogProvider,
	promocodeRepo repositories.Promocode,
) *handler {
	return &handler{
		b:               bot,
		customerService: customerService,
		orderService:    orderService,
		catalogProvider: catalogProvider,
		promocodeRepo:   promocodeRepo,
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
