package buttons

import (
	"domain"
	"household_bot/internal/telegram/callback"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	OrderTypeSelect = orderTypeSelectButtons()
)

func orderTypeSelectButtons() tg.InlineKeyboardMarkup {
	dataExpr := callback.Inject(callback.SelectOrderType, domain.OrderTypeExpress.String())
	dataNorm := callback.Inject(callback.SelectOrderType, domain.OrderTypeNormal.String())
	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("express", dataExpr),
		),

		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("normal", dataNorm),
		),
	)
}

func NewPaymentButton(c callback.Callback, orderShortID string) tg.InlineKeyboardMarkup {
	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Оплачено", callback.Inject(c, orderShortID)),
		))
}
