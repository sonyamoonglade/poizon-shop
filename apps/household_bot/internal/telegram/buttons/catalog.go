package buttons

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"household_bot/internal/telegram/callback"
)

func NewCategoryButtons(titles []string, cb callback.Callback, inStock bool, back BackButton) tg.InlineKeyboardMarkup {
	var rows [][]tg.InlineKeyboardButton
	for _, title := range titles {
		rows = append(rows, NewCategoryButton(cb, title, inStock).ToRow())
	}
	rows = append(rows, back.ToRow())
	return tg.NewInlineKeyboardMarkup(rows...)
}
func NewSubcategoryButtons(cTitle string, titles []string, cb callback.Callback, inStock bool, back BackButton) tg.InlineKeyboardMarkup {
	var rows [][]tg.InlineKeyboardButton
	for _, title := range titles {
		rows = append(rows, NewSubcategory(cb, cTitle, title, inStock).ToRow())
	}
	rows = append(rows, back.ToRow())
	return tg.NewInlineKeyboardMarkup(rows...)
}
