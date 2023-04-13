package handler

import (
	"context"
	"strconv"

	"domain"
	"dto"
	"household_bot/internal/telegram/templates"
	"household_bot/internal/telegram/tg_errors"
)

func (h *handler) AddToCart(ctx context.Context, chatID int64, args []string) error {
	var telegramID = chatID

	if err := h.checkRequiredState(ctx, telegramID, domain.StateWaitingToAddToCart); err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "AddToCart",
			CausedBy:    "checkRequiredState",
		})
	}
	if err := h.sendMessage(chatID, "Загружаем..."); err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "AddToCart",
			CausedBy:    "sendMessage",
		})
	}

	customer, err := h.customerRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "AddToCart",
			CausedBy:    "GetByTelegramID",
		})
	}
	var (
		cTitle,
		sTitle,
		inStockStr,
		pName = args[0], args[1], args[2], args[3]
	)
	inStock, err := strconv.ParseBool(inStockStr)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "AddToCart",
			CausedBy:    "ParseBool",
		})
	}
	products, err := h.categoryRepo.GetProductsByCategoryAndSubcategory(ctx, cTitle, sTitle, inStock)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "AddToCart",
			CausedBy:    "GetProductsByCategoryAndSubcategory",
		})
	}
	var p domain.HouseholdProduct
	for _, product := range products {
		if product.Name == pName {
			p = product
		}
	}

	customer.Cart.Add(p)

	err = h.customerRepo.Update(ctx, customer.CustomerID, dto.UpdateHouseholdCustomerDTO{
		Cart: &customer.Cart,
	})
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "AddToCart",
			CausedBy:    "Update",
		})
	}

	return h.sendMessage(chatID, templates.PositionAdded(pName))
}

func (h *handler) GetCart(ctx context.Context, chatID int64) error {
	var telegramID = chatID
	customer, err := h.customerRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "AddToCart",
			CausedBy:    "GetByTelegramID",
		})
	}
	if customer.Cart.IsEmpty() {
		return h.sendMessage(chatID, "empty cart")
	}
	return h.sendMessage(chatID, templates.RenderCart(templates.RenderCartArgs{
		Cart: customer.Cart,
	}))
}
