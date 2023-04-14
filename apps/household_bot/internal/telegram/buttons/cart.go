package buttons

import (
	"math"
	"strconv"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"household_bot/internal/telegram/callback"
)

var (
	CartPreview = cartPreviewButtons()
)

func cartPreviewButtons() tg.InlineKeyboardMarkup {
	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Оформить заказ", callback.Inject(callback.MakeOrder)),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Редактировать корзину", callback.Inject(callback.EditCart)),
			tg.NewInlineKeyboardButtonData("Добавить позицию", callback.Inject(callback.Catalog)),
		),
	)
}

func NewEditCartButtons(cartSize int, msgID int) tg.InlineKeyboardMarkup {
	keyboard := make([][]tg.InlineKeyboardButton, 0)

	var (
		numRows = int(math.Ceil(float64(cartSize) / 3))
		current int
	)
	for row := 0; row < numRows; row++ {
		keyboard = append(keyboard, tg.NewInlineKeyboardRow())
		for col := 0; col < 3 && current < cartSize; col++ {
			data := callback.Inject(callback.DeletePositionFromCart, strconv.Itoa(msgID), strconv.Itoa(current))
			button := tg.NewInlineKeyboardButtonData(strconv.Itoa(current+1), data)
			keyboard[row] = append(keyboard[row], button)
			current++
		}
	}

	return tg.NewInlineKeyboardMarkup(keyboard...)
}
