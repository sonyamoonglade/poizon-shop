package buttons

import (
	"strconv"

	"household_bot/internal/telegram/callback"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	RouteToCatalog = jumpToCatalog()
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

type ProductButtonsArgs struct {
	Cb             callback.Callback
	CTitle, STitle string
	Names          []string
	InStock        bool
	Back           BackButton
}

func NewProductsButtons(args ProductButtonsArgs) tg.InlineKeyboardMarkup {
	var rows [][]tg.InlineKeyboardButton

	for _, name := range args.Names {
		b := NewProductCard(args.Cb, args.CTitle, args.STitle, name, args.InStock)
		rows = append(rows, b.ToRow())
	}
	rows = append(rows, args.Back.ToRow())
	return tg.NewInlineKeyboardMarkup(rows...)
}

type ProductCardButtonsArgs struct {
	Cb                    callback.Callback
	CTitle, STitle, PName string
	InStock               bool
	Back                  BackButton
}

func NewProductCardButtons(args ProductCardButtonsArgs) tg.InlineKeyboardMarkup {
	rows := make([][]tg.InlineKeyboardButton, 0, 2)
	data := callback.Inject(args.Cb, args.CTitle, args.STitle, strconv.FormatBool(args.InStock), args.PName)
	addToCardBtn := tg.NewInlineKeyboardButtonData("Добавить в корзину", data)
	rows = append(rows, tg.NewInlineKeyboardRow(addToCardBtn))
	rows = append(rows, args.Back.ToRow())
	return tg.NewInlineKeyboardMarkup(rows...)
}

func jumpToCatalog() tg.InlineKeyboardMarkup {
	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("В каталог", callback.Inject(callback.Catalog)),
		),
	)
}
