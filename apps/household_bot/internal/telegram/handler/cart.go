package handler

import (
	"context"
	"fmt"
	"strconv"

	fn "github.com/elliotchance/pie/v2"

	"domain"
	"dto"
	"household_bot/internal/telegram/buttons"
	"household_bot/internal/telegram/templates"
	"household_bot/internal/telegram/tg_errors"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *handler) AddToCart(ctx context.Context, chatID int64, args []string) error {
	var telegramID = chatID

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

	products, err := h.categoryService.GetProductsByCategoryAndSubcategory(ctx, cTitle, sTitle, inStock)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "AddToCart",
			CausedBy:    "GetProductsByCategoryAndSubcategory",
		})
	}
	idx := fn.
		Of(products).
		FindFirstUsing(func(p domain.HouseholdProduct) bool {
			return p.Name == pName
		})
	currentProduct := products[idx]

	// Category that customer is adding product with
	currentCategory, err := h.categoryService.GetByID(ctx, currentProduct.CategoryID)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "AddToCart",
			CausedBy:    "GetByID",
		})
	}

	firstProduct, exists := customer.Cart.First()

	if customer.Cart.IsEmpty() {
		customer.Cart.Add(currentProduct)
	} else if exists {
		// Have to check if currentCategory is the same as 0th element in cart
		firstProductCategory, err := h.categoryService.GetByID(ctx, firstProduct.CategoryID)
		if err != nil {
			return tg_errors.New(tg_errors.Config{
				OriginalErr: err,
				Handler:     "AddToCart",
				CausedBy:    "GetByID",
			})
		}
		// If inStock field is the same then it's fine to add
		if firstProductCategory.InStock == currentCategory.InStock {
			customer.Cart.Add(currentProduct)
		} else {
			return h.sendWithKeyboard(
				chatID,
				templates.TryAddWithInvalidInStock(
					inStock,
					!inStock,
				),
				buttons.RouteToCatalogOrCart,
			)
		}
	}

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
	return h.sendWithKeyboard(chatID, templates.RenderCart(customer.Cart), buttons.CartPreview)
}

func (h *handler) EditCart(ctx context.Context, chatID int64, cartMsgID int) error {
	var telegramID = chatID

	customer, err := h.customerRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return fmt.Errorf("customerRepo.GetByTelegramID: %w", err)
	}

	if len(customer.Cart) == 0 {
		return h.emptyCart(chatID)
	}

	if err := h.customerRepo.UpdateState(ctx, telegramID, domain.StateWaitingForCartPositionToEdit); err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "EditCart",
			CausedBy:    "UpdateState",
		})
	}

	return h.sendWithKeyboard(chatID,
		templates.EditCartPosition(),
		buttons.NewEditCartButtons(
			len(customer.Cart),
			cartMsgID,
		),
	)
}

func (h *handler) DeletePositionFromCart(ctx context.Context, chatID int64, buttonsMsgID int, args []string) error {
	var (
		telegramID = chatID
		cartMsgIDStr,
		buttonClickedStr = args[0], args[1]
	)

	cartMsgID, err := strconv.Atoi(cartMsgIDStr)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "RemoveCartPosition",
			CausedBy:    "Atoi",
		})
	}

	posIndex, err := strconv.Atoi(buttonClickedStr)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "RemoveCartPosition",
			CausedBy:    "Atoi",
		})
	}

	if err := h.checkRequiredState(ctx, chatID, domain.StateWaitingForCartPositionToEdit); err != nil {
		return err
	}

	customer, err := h.customerRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "RemoveCartPosition",
			CausedBy:    "GetByTelegramID",
		})
	}

	customer.Cart.RemoveAt(posIndex)
	updateDTO := dto.UpdateHouseholdCustomerDTO{
		Cart: &customer.Cart,
	}
	if err := h.customerRepo.Update(ctx, customer.CustomerID, updateDTO); err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "RemoveCartPosition",
			CausedBy:    "Update",
		})
	}

	// If customer has emptied cart just now
	if customer.Cart.IsEmpty() {
		// Delete edit buttons
		if err := h.bot.CleanRequest(tg.NewDeleteMessage(chatID, buttonsMsgID)); err != nil {
			return tg_errors.New(tg_errors.Config{
				OriginalErr: err,
				Handler:     "RemoveCartPosition",
				CausedBy:    "CleanRequest",
			})
		}
		// Update cart message
		msg := tg.NewEditMessageText(chatID, cartMsgID, "Ваша корзина пуста!")
		keyboard := buttons.AddPosition
		msg.ReplyMarkup = &keyboard
		if err := h.cleanSend(msg); err != nil {
			return tg_errors.New(tg_errors.Config{
				OriginalErr: err,
				Handler:     "RemoveCartPosition",
				CausedBy:    "cleanSend",
			})
		}
		return nil
	}

	// Edit original preview cart message and edit buttons
	cartMsg := tg.NewEditMessageText(chatID, cartMsgID, templates.RenderCart(customer.Cart))
	cartMsg.ReplyMarkup = &buttons.CartPreview

	updateButtons := tg.NewEditMessageReplyMarkup(chatID,
		buttonsMsgID,
		buttons.NewEditCartButtons(
			len(customer.Cart),
			cartMsgID,
		),
	)

	if err := h.sendBulk(updateButtons, cartMsg); err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "RemoveCartPosition",
			CausedBy:    "sendBulk",
		})
	}

	return h.sendMessage(chatID, fmt.Sprintf("Позиция %d успешно удалена. Корзина сверху обновлена ✅", posIndex+1))
}

func (h *handler) emptyCart(chatID int64) error {
	return h.sendWithKeyboard(chatID, "Ваша корзина пуста!", buttons.AddPosition)
}
