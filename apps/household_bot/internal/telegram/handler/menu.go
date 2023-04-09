package handler

import (
	"context"

	"household_bot/internal/telegram/buttons"
)

func (h *handler) Start(ctx context.Context, chatID int64) error {
	return h.sendWithKeyboard(chatID, "start", buttons.Start)
}

func (h *handler) Menu(ctx context.Context, chatID int64) error {
	return h.sendWithKeyboard(chatID, "menu", buttons.Menu)
}

func (h *handler) Catalog(ctx context.Context, chatID int64) error {
	return h.sendWithKeyboard(chatID, "Выберите тип каталога", buttons.CatalogType)
}
