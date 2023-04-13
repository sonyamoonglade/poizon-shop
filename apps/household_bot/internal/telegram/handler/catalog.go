package handler

import (
	"context"
	"strconv"

	"domain"
	"functools"
	"household_bot/internal/telegram/buttons"
	"household_bot/internal/telegram/callback"
	"household_bot/internal/telegram/templates"
	"household_bot/internal/telegram/tg_errors"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	backButton := buttons.NewBackButton(callback.Catalog, nil, nil, &onlyAvailableInStock)
	editMsg := tg.NewEditMessageText(chatID, prevMsgID, "доступные категории")
	keyboard := buttons.NewCategoryButtons(categoryTitles, callback.SelectCategory, onlyAvailableInStock, backButton)
	editMsg.ReplyMarkup = &keyboard
	return h.cleanSend(editMsg)
}

func (h *handler) ProductsNew(ctx context.Context, chatID int64, msgIDForDeletion int, args []string) error {
	var (
		cTitle,
		sTitle,
		inStockStr = args[0], args[1], args[2]
	)
	inStock, err := strconv.ParseBool(inStockStr)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "ProductsNew",
			CausedBy:    "ParseBool",
		})
	}
	// p, ok := h.catalogProvider.GetProductAt(cTitle, sTitle, 0)
	// if !ok {
	// 	return tg_errors.New(tg_errors.Config{
	// 		OriginalErr: domain.ErrProductNotFound,
	// 		Handler:     "Products",
	// 		CausedBy:    "GetProductAt",
	// 	})
	// }

	products, err := h.householdCategoryRepo.GetProductsByCategoryAndSubcategory(ctx, cTitle, sTitle, inStock)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "ProductsNew",
			CausedBy:    "GetProductsByCategoryAndSubcategory",
		})
	}
	backButton := buttons.NewBackButton(callback.SelectCategory, &cTitle, &sTitle, &inStock)
	keyboard := buttons.NewProductsButtons(buttons.ProductButtonsArgs{
		CTitle:  cTitle,
		STitle:  sTitle,
		Cb:      callback.SelectProduct,
		InStock: inStock,
		Names: functools.Map(func(p domain.HouseholdProduct, _ int) string {
			return p.Name
		}, products),
		Back: backButton,
	})
	msg := tg.NewMessage(chatID, "Доутспные товары")
	msg.ReplyMarkup = keyboard
	// del card
	if err := h.bot.CleanRequest(tg.NewDeleteMessage(chatID, msgIDForDeletion)); err != nil {
		//todo: handler rr
		return err
	}
	return h.cleanSend(msg)

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
	backButton := buttons.NewBackButton(cb, &cTitle, nil, &onlyAvailableInStock)
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
	// p, ok := h.catalogProvider.GetProductAt(cTitle, sTitle, 0)
	// if !ok {
	// 	return tg_errors.New(tg_errors.Config{
	// 		OriginalErr: domain.ErrProductNotFound,
	// 		Handler:     "Products",
	// 		CausedBy:    "GetProductAt",
	// 	})
	// }

	products, err := h.householdCategoryRepo.GetProductsByCategoryAndSubcategory(ctx, cTitle, sTitle, onlyAvailableInStock)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "Products",
			CausedBy:    "GetProductsByCategoryAndSubcategory",
		})
	}
	backButton := buttons.NewBackButton(callback.SelectCategory, &cTitle, nil, &onlyAvailableInStock)
	keyboard := buttons.NewProductsButtons(buttons.ProductButtonsArgs{
		CTitle:  cTitle,
		STitle:  sTitle,
		Cb:      callback.SelectProduct,
		InStock: onlyAvailableInStock,
		Names: functools.Map(func(p domain.HouseholdProduct, _ int) string {
			return p.Name
		}, products),
		Back: backButton,
	})
	editMsg := tg.NewEditMessageText(chatID, prevMsgID, "Доутспные товары")
	editMsg.ReplyMarkup = &keyboard
	return h.cleanSend(editMsg)
}

func (h *handler) ProductCard(ctx context.Context, chatID int64, prevMsgID int, args []string) error {
	var (
		cTitle,
		sTitle,
		inStockStr,
		productName = args[0], args[1], args[2], args[3]
	)
	inStock, err := strconv.ParseBool(inStockStr)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "ProductCard",
			CausedBy:    "ParseBool",
		})
	}
	// Todo: replace with catalog provider
	products, err := h.householdCategoryRepo.GetProductsByCategoryAndSubcategory(ctx, cTitle, sTitle, inStock)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "ProductCard",
			CausedBy:    "GetProductsByCategoryAndSubcategory",
		})
	}

	var p domain.HouseholdProduct
	for _, product := range products {
		if product.Name == productName {
			p = product
		}
	}

	photo := tg.NewPhoto(chatID, tg.FileURL("https://picsum.photos/200/300"))
	photo.Caption = templates.HouseholdProductCaption(p)
	backButton := buttons.NewBackButton(callback.FromProductCardToProducts,
		&cTitle,
		&sTitle,
		&inStock)
	keyboard := tg.NewInlineKeyboardMarkup(backButton.ToRow())
	photo.ReplyMarkup = keyboard
	/// delete prev
	if err := h.bot.CleanRequest(tg.NewDeleteMessage(chatID, prevMsgID)); err != nil {
		// handle err
		return err
	}
	return h.cleanSend(photo)
}
