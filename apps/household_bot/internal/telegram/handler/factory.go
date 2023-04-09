package handler

import (
	"errors"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"household_bot/internal/catalog"
)

var (
	ErrInvalidState      = errors.New("invalid state")
	ErrInvalidPriceInput = errors.New("invalid price input")
)

type RateProvider interface {
	GetYuanRate() float64
}

type Bot interface {
	Send(c tg.Chattable) (tg.Message, error)
	CleanRequest(c tg.Chattable) error
	SendMediaGroup(c tg.MediaGroupConfig) ([]tg.Message, error)
}

type handler struct {
	bot             Bot
	rateProvider    RateProvider
	catalogProvider *catalog.Provider
}

func NewHandler(b Bot, rp RateProvider, cp *catalog.Provider) *handler {
	return &handler{
		bot:             b,
		rateProvider:    rp,
		catalogProvider: cp,
	}
}

func (h *handler) AnswerCallback(c *tg.CallbackQuery) error {
	return h.cleanSend(tg.NewCallback(c.ID, ""))
}
