package buttons

import (
	"strconv"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"household_bot/internal/telegram/callback"
)

// Following structs represent 'data' field in tg.InlineButton
type Category struct {
	c                    callback.Callback
	title                string
	onlyAvailableInStock bool
}

func NewCategoryButton(cb callback.Callback, title string, inStock bool) Category {
	return Category{
		c:                    cb,
		title:                title,
		onlyAvailableInStock: inStock,
	}
}

func (c Category) ToRow() []tg.InlineKeyboardButton {
	data := callback.Inject(c.c, c.title, strconv.FormatBool(c.onlyAvailableInStock))
	return tg.NewInlineKeyboardRow(tg.NewInlineKeyboardButtonData(c.title, data))
}

type Subcategory struct {
	c     callback.Callback
	title string
	// Parent category title
	cTitle               string
	onlyAvailableInStock bool
}

func NewSubcategory(cb callback.Callback, cTitle, title string, inStock bool) Subcategory {
	return Subcategory{
		c:                    cb,
		cTitle:               cTitle,
		title:                title,
		onlyAvailableInStock: inStock,
	}
}

func (s Subcategory) ToRow() []tg.InlineKeyboardButton {
	data := callback.Inject(s.c, s.cTitle, s.title, strconv.FormatBool(s.onlyAvailableInStock))
	return tg.NewInlineKeyboardRow(tg.NewInlineKeyboardButtonData(s.title, data))
}

type ProductCard struct {
	c              callback.Callback
	offset         int
	cTitle, sTitle string
}

func NewProductCard(cb callback.Callback, offset int, cTitle, sTitle string) ProductCard {
	return ProductCard{
		c:      cb,
		offset: offset,
		cTitle: cTitle,
		sTitle: sTitle,
	}
}

type BackButton struct {
	c       callback.Callback
	cTitle  *string
	inStock *bool
}

func NewBackButton(c callback.Callback, cTitle *string, inStock *bool) BackButton {
	return BackButton{
		c:       c,
		cTitle:  cTitle,
		inStock: inStock,
	}
}

const (
	backButtonTitle = "Вернуться назад"
)

func (b BackButton) ToRow() []tg.InlineKeyboardButton {
	var data string
	if b.cTitle != nil && b.inStock != nil {
		data = callback.Inject(b.c, *b.cTitle, strconv.FormatBool(*b.inStock))
	} else if b.inStock != nil {
		data = callback.Inject(b.c, strconv.FormatBool(*b.inStock))
	}

	return tg.NewInlineKeyboardRow(tg.NewInlineKeyboardButtonData(backButtonTitle, data))
}
