package telegram

import (
	"context"
	"fmt"

	"domain"
	"dto"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *handler) GetCart(ctx context.Context, chatID int64) error {
	var telegramID = chatID

	customer, err := h.customerRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return fmt.Errorf("customerRepo.GetByTelegramID: %w", err)
	}

	if len(customer.Cart) == 0 {
		return h.emptyCart(chatID)
	}
	isExpressOrder := *customer.Meta.NextOrderType == domain.OrderTypeExpress
	msg := tg.NewMessage(chatID, h.prepareCartPreview(customer.Cart, isExpressOrder))
	msg.ReplyMarkup = cartPreviewButtons

	return h.cleanSend(msg)
}

func (h *handler) EditCart(ctx context.Context, chatID int64, previewCartMsgID int) error {
	var telegramID = chatID

	customer, err := h.customerRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return fmt.Errorf("customerRepo.GetByTelegramID: %w", err)
	}

	if len(customer.Cart) == 0 {
		return h.emptyCart(chatID)
	}

	updateDTO := dto.UpdateClothingCustomerDTO{
		State: &domain.StateWaitingForCartPositionToEdit,
	}

	if err := h.customerRepo.Update(ctx, customer.CustomerID, updateDTO); err != nil {
		return fmt.Errorf("customerRepo.Update: %w", err)
	}
	buttons := prepareEditCartButtons(len(customer.Cart), previewCartMsgID)
	return h.sendWithKeyboard(chatID, editPositionTemplate, buttons)
}

func (h *handler) RemoveCartPosition(ctx context.Context, chatID int64, callbackData int, originalMsgID, cartPreviewMsgID int) error {
	var (
		telegramID    = chatID
		buttonClicked = callbackData - editCartRemovePositionOffset
		cartIndex     = buttonClicked - 1
	)

	if err := h.checkRequiredState(ctx, domain.StateWaitingForCartPositionToEdit, chatID); err != nil {
		return err
	}

	customer, err := h.customerRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return fmt.Errorf("customerRepo.GetByTelegramID: %w", err)
	}

	if buttonClicked > len(customer.Cart) {
		return fmt.Errorf("invalid button clicked")
	}

	customer.Cart.RemoveAt(cartIndex)

	updateDTO := dto.UpdateClothingCustomerDTO{
		Cart: &customer.Cart,
	}
	if err := h.customerRepo.Update(ctx, customer.CustomerID, updateDTO); err != nil {
		return fmt.Errorf("customerRepo.Update: %w", err)
	}

	// if customer has emptied cart
	if len(customer.Cart) == 0 {
		// delete edit buttons
		if err := h.b.CleanRequest(tg.NewDeleteMessage(chatID, originalMsgID)); err != nil {
			return fmt.Errorf("can't delete message: %w", err)
		}
		// update cartPreview
		msg := tg.NewEditMessageText(chatID, cartPreviewMsgID, "Ваша корзина пуста!")
		msg.ReplyMarkup = &addPositionButtons
		if err := h.cleanSend(msg); err != nil {
			return fmt.Errorf("cant edit cart preview message: %w", err)
		}
		return nil
	}

	// edit original preview cart message and edit buttons
	buttonsForNewCart := prepareEditCartButtons(len(customer.Cart), int(cartPreviewMsgID))

	isExpressOrder := *customer.Meta.NextOrderType == domain.OrderTypeExpress
	textForNewCart := h.prepareCartPreview(customer.Cart, isExpressOrder)

	updatePreviewText := tg.NewEditMessageText(chatID, int(cartPreviewMsgID), textForNewCart)
	updatePreviewText.ReplyMarkup = &cartPreviewButtons
	updateButtons := tg.NewEditMessageReplyMarkup(chatID, int(originalMsgID), buttonsForNewCart)

	if err := h.cleanSend(updateButtons); err != nil {
		return err
	}

	if err := h.cleanSend(updatePreviewText); err != nil {
		return err
	}

	return h.sendMessage(chatID, fmt.Sprintf("Позиция %d успешно удалена. Корзина сверху обновлена ✅", buttonClicked))
}

func (h *handler) emptyCart(chatID int64) error {
	return h.sendWithKeyboard(chatID, "Ваша корзина пуста!", addPositionButtons)
}

func (h *handler) prepareCartPreview(cart domain.ClothingCart, isExpressOrder bool) string {
	var out = getCartPreviewStartTemplate(len(cart), isExpressOrder)
	var totalRub uint64
	var totalYuan uint64
	for n, cartItem := range cart {
		positionText := getPositionTemplate(cartPositionPreviewArgs{
			n:         n + 1,
			link:      cartItem.ShopLink,
			size:      cartItem.Size,
			category:  string(cartItem.Category),
			priceRub:  cartItem.PriceRUB,
			priceYuan: cartItem.PriceYUAN,
		})
		totalRub += cartItem.PriceRUB
		totalYuan += cartItem.PriceYUAN
		out += positionText
	}
	out += getCartPreviewEndTemplate(totalRub, totalYuan)
	return out
}
