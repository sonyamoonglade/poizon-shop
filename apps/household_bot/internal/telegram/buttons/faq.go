package buttons

import (
	"strconv"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"household_bot/internal/telegram/callback"
	"household_bot/internal/telegram/templates"
)

var (
	FAQ     = faq()
	AskMore = askMore()
)

func faq() tg.InlineKeyboardMarkup {
	rows := make([][]tg.InlineKeyboardButton, 0, 8)
	for i := 0; i < 8; i++ {
		data := callback.Inject(callback.GetFaqAnswer, strconv.Itoa(i+1))
		button := tg.NewInlineKeyboardButtonData(templates.GetQuestion(i+1), data)
		rows = append(rows, tg.NewInlineKeyboardRow(button))
	}
	return tg.NewInlineKeyboardMarkup(rows...)
}

func askMore() tg.InlineKeyboardMarkup {
	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Жми", callback.Inject(callback.Faq)),
		))
}
