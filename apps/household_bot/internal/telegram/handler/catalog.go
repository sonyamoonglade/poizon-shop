package handler

import (
	"context"
	"fmt"
	"strconv"

	"domain"
	"functools"
	"household_bot/internal/telegram/buttons"
	"household_bot/internal/telegram/callback"
	"household_bot/internal/telegram/templates"
	"household_bot/internal/telegram/tg_errors"
	"household_bot/pkg/telegram"

	fn "github.com/elliotchance/pie/v2"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *handler) Categories(ctx context.Context, chatID int64, prevMsgID int, inStock bool) error {
	categories, err := h.categoryRepo.GetAllByInStock(ctx, inStock)
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

	onlyActive := fn.
		Of(categories).
		Filter(func(c domain.HouseholdCategory) bool {
			return c.Active
		}).
		Result
	categoryTitles := fn.Map(onlyActive, func(c domain.HouseholdCategory) string {
		return c.Title
	})

	// To prev step, reInject inStock
	backButton := buttons.NewBackButton(callback.Catalog, nil, nil, &inStock)
	text := fmt.Sprintf("Type: %s\nCategories", domain.InStockToString(inStock))
	editMsg := tg.NewEditMessageText(chatID, prevMsgID, text)
	keyboard := buttons.NewCategoryButtons(categoryTitles, callback.SelectCategory, inStock, backButton)
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
	// todo:into template
	text := fmt.Sprintf("Type: %s\nCategory: %s\nSubcategories:", domain.InStockToString(inStock), cTitle)
	editMsg := tg.NewEditMessageText(chatID, prevMsgID, text)
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
	if err := h.bot.CleanRequest(tg.NewDeleteMessage(chatID, prevMsgID)); err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "Products",
			CausedBy:    "CleanRequest",
		})
	}
	if err := h.catalogMsgService.DeleteByMsgID(ctx, prevMsgID); err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "Products",
			CausedBy:    "DeleteByMsgID",
		})
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

	photo := tg.NewPhoto(chatID, tg.FileURL(p.ImageURL))
	photo.Caption = templates.HouseholdProductCaption(p)
	photo.ParseMode = "markdown"
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

	return h.sendWithMessageID(photo, func(msgID int) error {
		catalogMsg := telegram.CatalogMsg{
			MsgID: msgID,
		}
		err := h.catalogMsgService.Save(ctx, catalogMsg)
		if err != nil {
			return tg_errors.New(tg_errors.Config{
				OriginalErr: err,
				Handler:     "ProductCard",
				CausedBy:    "Save",
			})
		}
		return nil
	})
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
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "ProductsNew",
			CausedBy:    "CleanRequest",
		})
	}
	if err := h.catalogMsgService.DeleteByMsgID(ctx, msgIDForDeletion); err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "ProductsNew",
			CausedBy:    "DeleteByMsgID",
		})
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

	return h.sendWithMessageID(c, func(msgID int) error {
		catalogMsg := telegram.CatalogMsg{
			MsgID: msgID,
		}
		err := h.catalogMsgService.Save(ctx, catalogMsg)
		if err != nil {
			return tg_errors.New(tg_errors.Config{
				OriginalErr: err,
				Handler:     "ProductsNew",
				CausedBy:    "Save",
			})
		}
		return nil
	})
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
	text := fmt.Sprintf("Type: %s\nCategory: %s\nSubcategory: %s\nProducts:", domain.InStockToString(inStock), cTitle, sTitle)
	if edit && editMsgID != nil {
		m := tg.NewEditMessageText(chatID, *editMsgID, text)
		m.ReplyMarkup = &keyboard
		c = m

	} else {
		m := tg.NewMessage(chatID, text)
		m.ReplyMarkup = keyboard
		c = m
	}
	return c, nil
}
