package handler

import (
	"context"
	"errors"

	"domain"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"household_bot/internal/telegram/buttons"
	"household_bot/internal/telegram/tg_errors"
	"household_bot/pkg/telegram"
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
	var c tg.Chattable
	if prevMsgID != nil {
		editMsg := tg.NewEditMessageText(chatID, *prevMsgID, "Выберите тип каталога")
		editMsg.ReplyMarkup = &buttons.CatalogType
		c = editMsg

	} else {
		msg := tg.NewMessage(chatID, "Выберите тип каталога")
		msg.ReplyMarkup = buttons.CatalogType
		c = msg
	}

	return h.sendWithMessageID(c, func(msgID int) error {
		catalogMsg := telegram.CatalogMsg{
			MsgID: msgID,
		}
		err := h.catalogMsgService.Save(ctx, catalogMsg)
		if err != nil {
			return tg_errors.New(tg_errors.Config{
				OriginalErr: err,
				Handler:     "Catalog",
				CausedBy:    "Save",
			})
		}
		return nil
	})
}
