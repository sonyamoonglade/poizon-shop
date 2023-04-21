package handler

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"domain"
	fn "github.com/sonyamoonglade/go_func"
	"household_bot/internal/telegram/buttons"
	"household_bot/internal/telegram/callback"
	"household_bot/internal/telegram/templates"
	"household_bot/internal/telegram/tg_errors"
	"household_bot/pkg/telegram"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *handler) Categories(ctx context.Context, chatID int64, prevMsgID int, inStock bool) error {
	categoryTitles := h.catalogProvider.GetActiveCategoryTitlesByInStock(inStock)
	if categoryTitles == nil {
		return h.sendMessage(chatID, "no categories")
	}

	// To prev step, reInject inStock
	backButton := buttons.NewBackButton(callback.Catalog, nil, nil, &inStock)
	text := fmt.Sprintf("Тип: %s\nКатегории:", domain.InStockToString(inStock))
	editMsg := tg.NewEditMessageText(chatID, prevMsgID, text)
	keyboard := buttons.NewCategoryButtons(categoryTitles, callback.SelectCategory, inStock, backButton)
	editMsg.ReplyMarkup = &keyboard
	return h.cleanSend(editMsg)
}

func (h *handler) Subcategories(ctx context.Context, chatID int64, prevMsgID int, args []string) error {
	cTitle := args[0]
	inStock, err := strconv.ParseBool(args[1])
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "Subcategories",
			CausedBy:    "ParseBool",
		})
	}

	subcategoryTitles := h.catalogProvider.GetSubcategoryTitles(cTitle, inStock)

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
	text := fmt.Sprintf("Тип: %s\nКатегория: %s\nПроизводители:", domain.InStockToString(inStock), cTitle)
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

	backButton := buttons.NewBackButton(callback.SelectCategory, &cTitle, nil, &inStock)
	c, err := h.fetchProductsAndGetChattable(
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

	keyboard := buttons.NewProductCardButtons(buttons.ProductCardButtonsArgs{
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
	customer, err := h.customerRepo.GetByTelegramID(ctx, chatID)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "ProductCard",
			CausedBy:    "GetByTelegramID",
		})
	}

	product := h.catalogProvider.GetProduct(cTitle, sTitle, productName, inStock)
	if err := h.renderProductCard(ctx, chatID, product, customer, keyboard); err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "ProductCard",
			CausedBy:    "renderProductCard",
		})
	}
	return nil
}

func (h *handler) ProductsNew(ctx context.Context, chatID int64, msgIDForDeletion int, args []string) error {
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
	if err := h.customerRepo.UpdateState(ctx, chatID, domain.StateDefault); err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "ProductsNew",
			CausedBy:    "UpdateState",
		})
	}

	backButton := buttons.NewBackButton(callback.SelectCategory, &cTitle, nil, &inStock)
	c, err := h.fetchProductsAndGetChattable(
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
			MsgID:  msgID,
			ChatID: chatID,
		}
		err := h.catalogMsgService.Save(ctx, catalogMsg)
		if err != nil {
			if errors.Is(err, telegram.ErrMessageAlreadyExists) {
				return nil
			}
			return tg_errors.New(tg_errors.Config{
				OriginalErr: err,
				Handler:     "ProductsNew",
				CausedBy:    "Save",
			})
		}
		return nil
	})
}

func (h *handler) fetchProductsAndGetChattable(
	chatID int64,
	edit bool,
	editMsgID *int,
	backButton buttons.BackButton,
	cTitle,
	sTitle string,
	inStock bool,
) (tg.Chattable, error) {
	products := h.catalogProvider.GetProducts(cTitle, sTitle, inStock)
	keyboard := buttons.NewProductsButtons(buttons.ProductButtonsArgs{
		CTitle:  cTitle,
		STitle:  sTitle,
		Cb:      callback.SelectProduct,
		InStock: inStock,
		Names: fn.Map(products, func(p domain.HouseholdProduct, _ int) string {
			return p.Name
		}),
		Back: backButton,
	})
	var c tg.Chattable
	text := fmt.Sprintf("Тип: %s\nКатегория: %s\nПроизводитель: %s\nМодели:", domain.InStockToString(inStock), cTitle, sTitle)
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

func (h *handler) renderProductCard(
	ctx context.Context,
	chatID int64,
	p domain.HouseholdProduct,
	customer domain.HouseholdCustomer,
	keyboard tg.InlineKeyboardMarkup,
) error {
	photo := tg.NewPhoto(chatID, tg.FileURL(p.ImageURL))
	if customer.HasPromocode() {
		promo, _ := customer.GetPromocode()
		photo.Caption = templates.HouseholdProductCaptionWithDiscount(p, promo.GetDiscount(domain.SourceHousehold))
	} else {
		photo.Caption = templates.HouseholdProductCaption(p)
	}
	photo.ParseMode = "markdown"
	photo.ReplyMarkup = keyboard
	return h.sendWithMessageID(photo, func(msgID int) error {
		catalogMsg := telegram.CatalogMsg{
			MsgID:  msgID,
			ChatID: chatID,
		}
		err := h.catalogMsgService.Save(ctx, catalogMsg)
		if err != nil {
			if errors.Is(err, telegram.ErrMessageAlreadyExists) {
				return nil
			}
			return tg_errors.New(tg_errors.Config{
				OriginalErr: err,
				Handler:     "renderProductCard",
				CausedBy:    "Save",
			})
		}
		return nil
	})
}
