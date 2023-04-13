package buttons

import (
	"strconv"

	"household_bot/internal/telegram/callback"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	c                    callback.Callback
	cTitle, sTitle, name string
	inStock              bool
}

func NewProductCard(cb callback.Callback, cTitle, sTitle, name string, inStock bool) ProductCard {
	return ProductCard{
		c:       cb,
		cTitle:  cTitle,
		sTitle:  sTitle,
		name:    name,
		inStock: inStock,
	}
}

func (p ProductCard) ToRow() []tg.InlineKeyboardButton {
	data := callback.Inject(p.c, p.cTitle, p.sTitle, strconv.FormatBool(p.inStock), p.name)
	return tg.NewInlineKeyboardRow(tg.NewInlineKeyboardButtonData(p.name, data))
}

type BackButton struct {
	c              callback.Callback
	cTitle, sTitle *string
	inStock        *bool
}

func NewBackButton(c callback.Callback, cTitle, sTitle *string, inStock *bool) BackButton {
	return BackButton{
		c:       c,
		cTitle:  cTitle,
		sTitle:  sTitle,
		inStock: inStock,
	}
}

const (
	backButtonTitle = "Вернуться назад"
)

func (b BackButton) ToRow() []tg.InlineKeyboardButton {
	var data string
	if b.cTitle != nil && b.sTitle != nil && b.inStock != nil {
		data = callback.Inject(b.c, *b.cTitle, *b.sTitle, strconv.FormatBool(*b.inStock))
	} else if b.cTitle != nil && b.inStock != nil {
		data = callback.Inject(b.c, *b.cTitle, strconv.FormatBool(*b.inStock))
	} else if b.inStock != nil {
		data = callback.Inject(b.c, strconv.FormatBool(*b.inStock))
	}
	return tg.NewInlineKeyboardRow(tg.NewInlineKeyboardButtonData(backButtonTitle, data))
}
