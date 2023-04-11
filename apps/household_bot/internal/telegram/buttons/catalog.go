package buttons

import (
	"household_bot/internal/telegram/callback"
	"strconv"

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

type ProductCarouselButtonsArgs struct {
	CTitle, STitle       string
	InStock              bool
	CurrOffset           int
	HasNext, HasPrev     bool
	NextTitle, PrevTitle *string
	Back                 BackButton
}

const (
	arrLeft, arrRight = "<", ">"
)

func NewProductCarouselButtons(args ProductCarouselButtonsArgs) tg.InlineKeyboardMarkup {
	var rows [][]tg.InlineKeyboardButton
	if args.HasNext && args.NextTitle != nil && !args.HasPrev {
		data := callback.Inject(callback.NextProduct,
			args.CTitle,
			args.STitle,
			strconv.FormatBool(args.InStock),
			strconv.Itoa(args.CurrOffset+1),
		)
		nextBtn := tg.NewInlineKeyboardButtonData(*args.NextTitle, data)
		rows = append(rows, []tg.InlineKeyboardButton{nextBtn})
	}
	if args.HasPrev && !args.HasNext && args.PrevTitle != nil {
		data := callback.Inject(callback.PrevProduct,
			args.CTitle,
			args.STitle,
			strconv.FormatBool(args.InStock),
			strconv.Itoa(args.CurrOffset-1),
		)
		prevBtn := tg.NewInlineKeyboardButtonData(*args.PrevTitle, data)
		rows = append(rows, []tg.InlineKeyboardButton{prevBtn})
	}
	if args.HasNext && args.HasPrev && args.PrevTitle != nil && args.NextTitle != nil {
		dataNext := callback.Inject(callback.NextProduct,
			args.CTitle,
			args.STitle,
			strconv.FormatBool(args.InStock),
			strconv.Itoa(args.CurrOffset+1),
		)
		dataPrev := callback.Inject(callback.PrevProduct,
			args.CTitle,
			args.STitle,
			strconv.FormatBool(args.InStock),
			strconv.Itoa(args.CurrOffset-1),
		)
		prevBtn := tg.NewInlineKeyboardButtonData(*args.PrevTitle, dataPrev)
		nextBtn := tg.NewInlineKeyboardButtonData(*args.NextTitle, dataNext)
		rows = append(rows, []tg.InlineKeyboardButton{prevBtn, nextBtn})
	}

	rows = append(rows, args.Back.ToRow())

	return tg.NewInlineKeyboardMarkup(rows...)
}
