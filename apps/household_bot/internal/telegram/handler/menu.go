package handler

import (
	"context"
	"errors"

	"domain"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"household_bot/internal/telegram/buttons"
	"household_bot/internal/telegram/tg_errors"
)

func (h *handler) Start(ctx context.Context, m *tg.Message) error {
	var telegramID = m.Chat.ID
	_, err := h.customerRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		if !errors.Is(err, domain.ErrCustomerNotFound) {
			return tg_errors.New(tg_errors.Config{
				OriginalErr: err,
				Handler:     "Start",
				CausedBy:    "GetByTelegramID",
			})

		}

		err := h.customerRepo.Save(ctx, domain.NewHouseholdCustomer(telegramID, domain.MakeUsername(m.From.String())))
		if err != nil {
			return tg_errors.New(tg_errors.Config{
				OriginalErr: err,
				Handler:     "Start",
				CausedBy:    "Save",
			})
		}
	}
	return h.sendWithKeyboard(telegramID, "start", buttons.Start)
}

func (h *handler) Menu(ctx context.Context, chatID int64) error {
	// cleanup catalogMsgIDs...
	if err := h.customerRepo.UpdateState(ctx, chatID, domain.StateDefault); err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "Menu",
			CausedBy:    "UpdateState",
		})
	}
	return h.sendWithKeyboard(chatID, "menu", buttons.Menu)
}

func (h *handler) Catalog(ctx context.Context, chatID int64, prevMsgID *int) error {
	if prevMsgID != nil {
		editMsg := tg.NewEditMessageText(chatID, *prevMsgID, "Выберите тип каталога")
		editMsg.ReplyMarkup = &buttons.CatalogType
		return h.cleanSend(editMsg)
	}
	return h.sendWithKeyboard(chatID, "Выберите тип каталога", buttons.CatalogType)
}
