package handler

import (
	"context"
	"strconv"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"household_bot/internal/telegram/buttons"
	"household_bot/internal/telegram/callback"
	"household_bot/internal/telegram/tg_errors"
)

func (h *handler) Categories(ctx context.Context, chatID int64, prevMsgID int, onlyAvailableInStock bool) error {
	categories, err := h.householdCategoryRepo.GetAll(ctx)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "Categories",
			CausedBy:    "GetAll",
		})
	}
	if categories == nil {
		return h.sendMessage(chatID, "no categories")
	}
	var categoryTitles []string
	for _, c := range categories {
		if c.Active {
			categoryTitles = append(categoryTitles, c.Title)
		}
	}

	// To prev step, reInject inStock
	backButton := buttons.NewBackButton(callback.Catalog, nil, &onlyAvailableInStock)
	editMsg := tg.NewEditMessageText(chatID, prevMsgID, "доступные категории")
	keyboard := buttons.NewCategoryButtons(categoryTitles, callback.SelectCategory, onlyAvailableInStock, backButton)
	editMsg.ReplyMarkup = &keyboard
	return h.cleanSend(editMsg)
}

func (h *handler) Subcategories(ctx context.Context, chatID int64, prevMsgID int, args []string) error {
	onlyAvailableInStock, err := strconv.ParseBool(args[1])
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "Subcategories",
			CausedBy:    "ParseBool",
		})
	}

	cTitle := args[0]
	category, err := h.householdCategoryRepo.GetByTitle(ctx, cTitle)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "Subcategories",
			CausedBy:    "GetByTitle",
		})
	}

	var subcategoryTitles []string
	for _, s := range category.Subcategories {
		if s.Active {
			subcategoryTitles = append(subcategoryTitles, s.Title)
		}
	}
	cb := callback.CTypeOrder
	if onlyAvailableInStock {
		cb = callback.CTypeInStock
	}
	// To prev step, reInject inStock and cTitle
	backButton := buttons.NewBackButton(cb, &cTitle, &onlyAvailableInStock)
	keyboard := buttons.NewSubcategoryButtons(cTitle, subcategoryTitles, callback.SelectSubcategory, onlyAvailableInStock, backButton)
	editMsg := tg.NewEditMessageText(chatID, prevMsgID, "Подкатегории")
	editMsg.ReplyMarkup = &keyboard
	return h.cleanSend(editMsg)
}

func (h *handler) Products(ctx context.Context, chatID int64, prevMsgID int, args []string) error {
	onlyAvailableInStock, err := strconv.ParseBool(args[2])
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "Products",
			CausedBy:    "ParseBool",
		})
	}
	cTitle := args[0]
	sTitle := args[1]
	products, err := h.householdCategoryRepo.GetProductsByCategoryAndSubcategory(ctx, cTitle, sTitle, onlyAvailableInStock)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "Products",
			CausedBy:    "GetProductsByCategoryAndSubcategory",
		})
	}
	//todo: edit msg, render nicely as scroll
	for _, p := range products {
		h.sendMessage(chatID, p.Name+" "+p.ISBN)
	}
	return nil
}
