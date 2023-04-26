package telegram

import (
	"context"
	"fmt"

	"domain"
	"dto"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *handler) GetCart(ctx context.Context, chatID int64) error {
	customer, err := h.customerService.GetByTelegramID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("customerRepo.GetByTelegramID: %w", err)
	}
	if customer.Cart.Empty() {
		return h.emptyCart(chatID)
	}

	isExpressOrder := customer.Meta.
		NextOrderType.
		IsExpress()

	msg := tg.NewMessage(
		chatID,
		h.prepareCartPreview(
			customer.Cart,
			isExpressOrder,
			customer.HasPromocode(),
			customer.MustGetPromocode().GetClothingDiscount(),
		),
	)
	msg.ReplyMarkup = cartPreviewButtons
	return h.cleanSend(msg)
}

func (h *handler) EditCart(ctx context.Context, chatID int64, previewCartMsgID int) error {
	var telegramID = chatID

	customer, err := h.customerService.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return fmt.Errorf("customerRepo.GetByTelegramID: %w", err)
	}

	if len(customer.Cart) == 0 {
		return h.emptyCart(chatID)
	}

	updateDTO := dto.UpdateClothingCustomerDTO{
		State: &domain.StateWaitingForCartPositionToEdit,
	}

	if err := h.customerService.Update(ctx, customer.CustomerID, updateDTO); err != nil {
		return fmt.Errorf("customerRepo.Update: %w", err)
	}
	buttons := prepareEditCartButtons(len(customer.Cart), previewCartMsgID)
	return h.sendWithKeyboard(chatID, editPositionTemplate, buttons)
}

func (h *handler) RemoveCartPosition(ctx context.Context, chatID int64, callbackData int, originalMsgID, cartPreviewMsgID int) error {
	var (
		buttonClicked = callbackData - editCartRemovePositionOffset
		cartIndex     = buttonClicked - 1
	)

	customer, err := h.checkRequiredState(ctx, chatID, domain.StateWaitingForCartPositionToEdit)
	if err != nil {
		return fmt.Errorf("check required state: %w", err)
	}

	if buttonClicked > len(customer.Cart) {
		return fmt.Errorf("invalid button clicked")
	}

	customer.Cart.RemoveAt(cartIndex)

	updateDTO := dto.UpdateClothingCustomerDTO{
		Cart: &customer.Cart,
	}
	if err := h.customerService.Update(ctx, customer.CustomerID, updateDTO); err != nil {
		return fmt.Errorf("customerRepo.Update: %w", err)
	}

	// if customer has emptied cart
	if customer.Cart.Empty() {
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
	buttonsForNewCart := prepareEditCartButtons(len(customer.Cart), cartPreviewMsgID)

	isExpressOrder := customer.Meta.NextOrderType.IsExpress()
	textForNewCart := h.prepareCartPreview(
		customer.Cart,
		isExpressOrder,
		customer.HasPromocode(),
		customer.MustGetPromocode().GetClothingDiscount(),
	)

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

func (h *handler) prepareCartPreview(cart domain.ClothingCart, isExpressOrder bool, discounted bool, discount uint32) string {
	var out = getCartPreviewStartTemplate(len(cart), isExpressOrder)
	var totalRub uint64
	var totalYuan uint64
	for n, cartItem := range cart {
		if discounted {
			out += getDiscountedPositionTemplate(cartPositionPreviewDiscountedArgs{
				n:           n + 1,
				link:        cartItem.ShopLink,
				size:        cartItem.Size,
				discountRub: discount,
				category:    string(cartItem.Category),
				priceRub:    cartItem.PriceRUB,
				priceYuan:   cartItem.PriceYUAN,
			})
			totalRub += cartItem.PriceRUB - uint64(discount)
			totalYuan += cartItem.PriceYUAN
			continue
		}
		out += getPositionTemplate(cartPositionPreviewArgs{
			n:         n + 1,
			link:      cartItem.ShopLink,
			size:      cartItem.Size,
			category:  string(cartItem.Category),
			priceRub:  cartItem.PriceRUB,
			priceYuan: cartItem.PriceYUAN,
		})
		totalRub += cartItem.PriceRUB
		totalYuan += cartItem.PriceYUAN
	}
	out += getCartPreviewEndTemplate(totalRub, totalYuan)
	return out
}
