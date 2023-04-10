package handler

import (
	"context"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"household_bot/internal/telegram/buttons"
)

func (h *handler) Start(ctx context.Context, chatID int64) error {
	return h.sendWithKeyboard(chatID, "start", buttons.Start)
}

func (h *handler) Menu(ctx context.Context, chatID int64) error {
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
