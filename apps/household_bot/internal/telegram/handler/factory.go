package handler

import (
	"context"
	"errors"
	"household_bot/internal/catalog"

	"repositories"

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
	bot                   Bot
	rateProvider          RateProvider
	householdCategoryRepo repositories.HouseholdCategory
	catalogProvider       *catalog.Provider
}

func NewHandler(b Bot, rp RateProvider, householdCategoryRepo repositories.HouseholdCategory, catalogProvider *catalog.Provider) *handler {
	return &handler{
		bot:                   b,
		rateProvider:          rp,
		householdCategoryRepo: householdCategoryRepo,
		catalogProvider:       catalogProvider,
	}
}

func (h *handler) AnswerCallback(c *tg.CallbackQuery) error {
	return h.cleanSend(tg.NewCallback(c.ID, ""))
}
