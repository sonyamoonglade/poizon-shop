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
	categories, err := h.categoryRepo.GetAll(ctx)
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

func (h *handler) Subcategories(ctx context.Context, chatID int64, prevMsgID int, args []string) error {
	inStock, err := strconv.ParseBool(args[1])
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "Subcategories",
			CausedBy:    "ParseBool",
		})
	}

	cTitle := args[0]
	category, err := h.categoryRepo.GetByTitle(ctx, cTitle, inStock)
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
	if inStock {
		cb = callback.CTypeInStock
	}
	// To prev step, reInject inStock and cTitle
	backButton := buttons.NewBackButton(cb, &cTitle, nil, &inStock)
	keyboard := buttons.NewSubcategoryButtons(
		cTitle,
		subcategoryTitles,
		callback.SelectSubcategory,
		inStock,
		backButton,
	)
	editMsg := tg.NewEditMessageText(chatID, prevMsgID, "Подкатегории")
	editMsg.ReplyMarkup = &keyboard
	return h.cleanSend(editMsg)
}

func (h *handler) Products(ctx context.Context, chatID int64, prevMsgID int, args []string) error {
	var (
		cTitle,
		sTitle,
		inStockStr = args[0], args[1], args[2]
	)
	inStock, err := strconv.ParseBool(inStockStr)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "Products",
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
	backButton := buttons.NewBackButton(callback.SelectCategory, &cTitle, nil, &inStock)
	c, err := h.fetchProductsAndGetChattable(
		ctx,
		chatID,
		true,
		&prevMsgID,
		backButton,
		cTitle,
		sTitle,
		inStock,
	)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "Products",
			CausedBy:    "fetchProductsAndGetChattable",
		})
	}
	return h.cleanSend(c)
}

func (h *handler) ProductCard(ctx context.Context, chatID int64, prevMsgID int, args []string) error {
	/// delete prev
	if err := h.bot.CleanRequest(tg.NewDeleteMessage(chatID, prevMsgID)); err != nil {
		// handle err
		return err
	}
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
	products, err := h.categoryRepo.GetProductsByCategoryAndSubcategory(ctx, cTitle, sTitle, inStock)
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
	photo.ReplyMarkup = buttons.NewProductCardButtons(buttons.ProductCardButtonsArgs{
		Cb:      callback.AddToCart,
		CTitle:  cTitle,
		STitle:  sTitle,
		PName:   productName,
		InStock: inStock,
		Back: buttons.NewBackButton(callback.FromProductCardToProducts,
			&cTitle,
			&sTitle,
			&inStock,
		),
	})

	if err := h.customerRepo.UpdateState(ctx, chatID, domain.StateWaitingToAddToCart); err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "ProductCard",
			CausedBy:    "UpdateState",
		})
	}

	return h.cleanSend(photo)
}

func (h *handler) ProductsNew(ctx context.Context,
	chatID int64,
	msgIDForDeletion int,
	args []string,
) error {
	var (
		cTitle,
		sTitle,
		inStockStr = args[0], args[1], args[2]
	)
	if err := h.bot.CleanRequest(tg.NewDeleteMessage(chatID, msgIDForDeletion)); err != nil {
		//todo: handler rr
		return err
	}
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
	if err := h.customerRepo.UpdateState(ctx, chatID, domain.StateDefault); err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "ProductsNew",
			CausedBy:    "UpdateState",
		})
	}

	backButton := buttons.NewBackButton(callback.SelectCategory, &cTitle, nil, &inStock)
	c, err := h.fetchProductsAndGetChattable(
		ctx,
		chatID,
		false,
		nil,
		backButton,
		cTitle,
		sTitle,
		inStock,
	)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "ProductsNew",
			CausedBy:    "fetchProductsAndGetChattable",
		})
	}
	return h.cleanSend(c)
}

func (h *handler) fetchProductsAndGetChattable(ctx context.Context,
	chatID int64,
	edit bool,
	editMsgID *int,
	backButton buttons.BackButton,
	cTitle,
	sTitle string,
	inStock bool,
) (tg.Chattable, error) {
	products, err := h.categoryRepo.GetProductsByCategoryAndSubcategory(ctx, cTitle, sTitle, inStock)
	if err != nil {
		return nil, tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "fetchProductsAndGetChattable",
			CausedBy:    "GetProductsByCategoryAndSubcategory",
		})
	}
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
	var c tg.Chattable
	if edit && editMsgID != nil {
		m := tg.NewEditMessageText(chatID, *editMsgID, "Доступные товары")
		m.ReplyMarkup = &keyboard
		c = m

	} else {
		m := tg.NewMessage(chatID, "Доступные товары")
		m.ReplyMarkup = keyboard
		c = m
	}
	return c, nil
}
