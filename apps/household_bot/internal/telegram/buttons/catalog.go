package buttons

import (
	"household_bot/internal/telegram/callback"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
