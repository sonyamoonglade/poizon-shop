package buttons

import (
	"strconv"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"household_bot/internal/telegram/callback"
)

// Following structs represent 'data' field in tg.InlineButton
type Category struct {
	c     callback.Callback
	msgID int
	title string
}

func NewCategoryButton(cb callback.Callback, msgID int, title string) Category {
	return Category{
		c:     cb,
		msgID: msgID,
		title: title,
	}
}

func (c Category) ToRow() []tg.InlineKeyboardButton {
	return tg.NewInlineKeyboardRow(tg.NewInlineKeyboardButtonData(c.title, callback.Inject(c.c, strconv.Itoa(c.msgID))))
}

type Subcategory struct {
	c     callback.Callback
	msgID int
	title string
	// Parent category title
	cTitle string
}

func NewSubcategory(cb callback.Callback, msgID int, cTitle, title string) Subcategory {
	return Subcategory{
		c:      cb,
		msgID:  msgID,
		cTitle: cTitle,
		title:  title,
	}
}

type ProductCard struct {
	c              callback.Callback
	offset         int
	msgID          int
	cTitle, sTitle string
}

func NewProductCart(cb callback.Callback, msgID, offset int, cTitle, sTitle string) ProductCard {
	return ProductCard{
		c:      cb,
		msgID:  msgID,
		offset: offset,
		cTitle: cTitle,
		sTitle: sTitle,
	}
}
