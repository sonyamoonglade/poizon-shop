package handler

import (
	"context"

	"household_bot/internal/telegram/buttons"
	"household_bot/internal/telegram/callback"
)

func (h *handler) Categories(ctx context.Context, chatID int64, prevMsgID int, onlyAvailableInStock bool) error {
	categoryTitles := h.catalogProvider.GetCategoryTitles(true)
	return h.sendWithKeyboard(chatID, "Категории", buttons.NewCategoryButtons(categoryTitles, callback.SelectCategory, prevMsgID))
}
