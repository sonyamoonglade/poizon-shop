package buttons

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"household_bot/internal/telegram/callback"
)

func NewCategoryButtons(titles []string, cb callback.Callback, msgID int) tg.InlineKeyboardMarkup {
	var rows [][]tg.InlineKeyboardButton
	for _, title := range titles {
		rows = append(rows, NewCategoryButton(cb, msgID, title).ToRow())
	}
	return tg.NewInlineKeyboardMarkup(rows...)
}
