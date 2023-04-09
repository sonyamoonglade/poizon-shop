package buttons

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"household_bot/internal/telegram/callback"
	"household_bot/internal/telegram/router"
)

var (
	Menu        = menu()
	Start       = start()
	CatalogType = catalogType()
)

func menu() tg.InlineKeyboardMarkup {
	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Каталог", callback.Inject(callback.Catalog)),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Вопросы", callback.Inject(callback.Faq)),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Мои заказы", callback.Inject(callback.MyOrders)),
		),
	)
}

func start() tg.ReplyKeyboardMarkup {
	return tg.NewReplyKeyboard(
		tg.NewKeyboardButtonRow(
			tg.NewKeyboardButton(router.Menu),
		),
	)
}

func catalogType() tg.InlineKeyboardMarkup {
	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("В наличии", callback.Inject(callback.CTypeInStock)),
			tg.NewInlineKeyboardButtonData("Под заказ", callback.Inject(callback.CTypeOrder)),
		),
	)
}
