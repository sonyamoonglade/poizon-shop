package handler

import (
	"context"
	"errors"
	"fmt"

	"domain"
	"household_bot/internal/catalog"
	"household_bot/pkg/telegram"
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
	catalogProvider *catalog.Provider
	promocodeRepo   repositories.Promocode

	categoryService   services.HouseholdCategory
	orderService      services.Order[domain.HouseholdOrder]
	catalogMsgService services.HouseholdCatalogMsg
	customerService   services.HouseholdCustomer
}

func NewHandler(b Bot,
	rp RateProvider,
	repos repositories.Repositories,
	catalogProvider *catalog.Provider,
	orderService services.Order[domain.HouseholdOrder],
	catalogMsgService services.HouseholdCatalogMsg,
	categoryService services.HouseholdCategory,
	customerService services.HouseholdCustomer,
) *handler {
	return &handler{
		bot:               b,
		rateProvider:      rp,
		categoryService:   categoryService,
		promocodeRepo:     repos.Promocode,
		catalogProvider:   catalogProvider,
		orderService:      orderService,
		catalogMsgService: catalogMsgService,
		customerService:   customerService,
	}
}

func (h *handler) AnswerCallback(c *tg.CallbackQuery) error {
	return h.cleanSend(tg.NewCallback(c.ID, ""))
}

const sorry = "Упс, что-то пошло не так."

func (h *handler) Sorry(chatID int64) error {
	return h.sendMessage(chatID, sorry)
}

// WipeCatalogs reads all msgID's of catalogs from database and deletes them
func (h *handler) WipeCatalogs(ctx context.Context) error {
	err := h.catalogMsgService.WipeAll(ctx, h)
	if err != nil {
		return fmt.Errorf("wipe catalogs: %w", err)
	}
	return nil
}

func (h *handler) DeleteFromCatalog(m telegram.CatalogMsg) error {
	return h.bot.CleanRequest(tg.NewDeleteMessage(m.ChatID, m.MsgID))
}
