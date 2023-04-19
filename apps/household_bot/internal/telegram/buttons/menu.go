package buttons

import (
	"household_bot/internal/telegram/callback"
	"household_bot/internal/telegram/router"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	Menu        = menu()
	Start       = start()
	CatalogType = catalogType()
	AddPosition = addPosition()
	MakeOrder   = makeOrder()
)

func menu() tg.InlineKeyboardMarkup {
	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Каталог", callback.Inject(callback.Catalog)),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Search by ISBN", callback.Inject(callback.GetProductByISBN)),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Вопросы", callback.Inject(callback.Faq)),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Мои заказы", callback.Inject(callback.MyOrders)),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Посмотреть корзину", callback.Inject(callback.MyCart)),
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

func addPosition() tg.InlineKeyboardMarkup {
	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Добавить позицию", callback.Inject(callback.Catalog)),
		))
}

func makeOrder() tg.InlineKeyboardMarkup {
	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Сделать заказ", callback.Inject(callback.Catalog)),
		))
}
