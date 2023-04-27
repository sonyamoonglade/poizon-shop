package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"domain"
	"dto"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"utils/url"
)

func (h *handler) AddPosition(ctx context.Context, chatID int64) error {
	var (
		telegramID = chatID
	)
	if err := h.sendMessage(chatID, newPositionWarnTemplate); err != nil {
		return err
	}
	customer, err := h.customerService.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return err
	}
	if len(customer.Cart) == 0 {
		if err := h.customerService.UpdateState(ctx, telegramID, domain.StateWaitingForOrderType); err != nil {
			return err
		}
		// Start from scratch
		return h.askForOrderType(ctx, chatID)
	}
	if err := h.customerService.UpdateState(ctx, telegramID, domain.StateWaitingForCategory); err != nil {
		return err
	}
	// Otherwise start from category selection
	return h.askForCategory(ctx, chatID)
}

func (h *handler) askForCategory(ctx context.Context, chatID int64) error {
	if err := h.customerService.UpdateState(ctx, chatID, domain.StateWaitingForCategory); err != nil {
		return err
	}
	return h.sendWithKeyboard(chatID, askForCategoryTemplate, categoryButtons)
}

func (h *handler) HandleCategoryInput(ctx context.Context, chatID int64, cat domain.Category) error {
	customer, err := h.checkRequiredState(ctx, chatID, domain.StateWaitingForCategory)
	if err != nil {
		return err
	}

	customer.UpdateLastEditPositionCategory(cat)
	updateDTO := dto.UpdateClothingCustomerDTO{
		LastPosition: customer.LastEditPosition,
	}

	if err := h.customerService.Update(ctx, customer.CustomerID, updateDTO); err != nil {
		return fmt.Errorf("customerRepo.Update: %w", err)
	}

	if err := h.sendMessage(chatID, fmt.Sprintf("Выбрана категория: %s", string(cat))); err != nil {
		return err
	}

	return h.askForSize(ctx, chatID)
}

func (h *handler) askForOrderType(ctx context.Context, chatID int64) error {
	text := "Выбери тип доставки"
	return h.sendWithKeyboard(chatID, text, orderTypeButtons)
}

func (h *handler) HandleOrderTypeInput(ctx context.Context, chatID int64, typ domain.OrderType) error {
	var (
		isExpress = typ == domain.OrderTypeExpress
	)

	customer, err := h.checkRequiredState(ctx, chatID, domain.StateWaitingForOrderType)
	if err != nil {
		return fmt.Errorf("check required state: %w", err)
	}

	customer.UpdateMetaOrderType(typ)

	updateDTO := dto.UpdateClothingCustomerDTO{
		Meta:  &customer.Meta,
		State: &domain.StateWaitingForCategory,
	}
	if err := h.customerService.Update(ctx, customer.CustomerID, updateDTO); err != nil {
		return err
	}

	var resp = "Тип доставки: "
	switch isExpress {
	case true:
		resp += "Экспресс"
		break
	case false:
		resp += "Обычный"
		break
	}
	if err := h.sendMessage(chatID, resp); err != nil {
		return err
	}

	return h.askForCategory(ctx, chatID)
}

func (h *handler) askForSize(ctx context.Context, chatID int64) error {
	if err := h.sendWithKeyboard(chatID, askForSizeTemplate, bottomMenuWithoutAddPositionButtons); err != nil {
		return err
	}

	return h.customerService.UpdateState(ctx, chatID, domain.StateWaitingForSize)
}

func (h *handler) HandleSizeInput(ctx context.Context, m *tg.Message) error {
	var (
		chatID   = m.Chat.ID
		sizeText = strings.TrimSpace(m.Text)
	)

	customer, err := h.checkRequiredState(ctx, chatID, domain.StateWaitingForSize)
	if err != nil {
		return fmt.Errorf("check required state: %w", err)
	}

	customer.UpdateLastEditPositionSize(sizeText)

	updateDTO := dto.UpdateClothingCustomerDTO{
		LastPosition: customer.LastEditPosition,
		State:        &domain.StateWaitingForButton,
	}

	if err := h.customerService.Update(ctx, customer.CustomerID, updateDTO); err != nil {
		return fmt.Errorf("customerRepo.Update: %w", err)
	}
	if sizeText == "#" {
		sizeText = "БЕЗ размера"
	}
	if err := h.sendMessage(chatID, fmt.Sprintf("Твой размер: %s", sizeText)); err != nil {
		return err
	}
	return h.sendWithKeyboard(chatID, askForButtonColorTemplate, selectColorButtons)
}

func (h *handler) HandleButtonSelect(ctx context.Context, c *tg.CallbackQuery, button domain.Button) error {
	var (
		chatID = c.From.ID
	)

	customer, err := h.checkRequiredState(ctx, chatID, domain.StateWaitingForButton)
	if err != nil {
		return fmt.Errorf("check required state: %w", err)
	}

	customer.UpdateLastEditPositionButtonColor(button)
	updateDTO := dto.UpdateClothingCustomerDTO{
		LastPosition: customer.LastEditPosition,
		State:        &domain.StateWaitingForPrice,
	}

	if err := h.customerService.Update(ctx, customer.CustomerID, updateDTO); err != nil {
		return fmt.Errorf("customerRepo.Update: %w", err)
	}
	if err := h.sendMessage(chatID, fmt.Sprintf("Цвет выбранной кнопки: %s", string(button))); err != nil {
		return err
	}

	return h.sendMessage(chatID, askForPriceTemplate)
}

func (h *handler) HandlePriceInput(ctx context.Context, m *tg.Message) error {
	var (
		chatID = m.Chat.ID
		input  = strings.TrimSpace(m.Text)
	)
	// validate state
	customer, err := h.checkRequiredState(ctx, chatID, domain.StateWaitingForPrice)
	if err != nil {
		return fmt.Errorf("check required state: %w", err)
	}

	priceYuan, err := strconv.ParseUint(input, 10, 64)
	if err != nil {
		return ErrInvalidPriceInput
	}
	var (
		ordTyp = customer.Meta.NextOrderType
	)
	if ordTyp == nil {
		return fmt.Errorf("order type in meta is nil")
	}
	// We should apply customer.Meta and customer.LastEditPosition.Category in order to calculate correctly
	rate, err := h.rateProvider.GetYuanRate(ctx)
	if err != nil {
		return fmt.Errorf("get yuan rate: %w", err)
	}
	args := domain.ConvertYuanArgs{
		X:         priceYuan,
		Rate:      rate,
		OrderType: *ordTyp,
		Category:  customer.LastEditPosition.Category,
	}

	priceRub := domain.ConvertYuan(args)
	customer.UpdateLastEditPositionPrice(priceRub, priceYuan)

	updateDTO := dto.UpdateClothingCustomerDTO{
		LastPosition: customer.LastEditPosition,
		State:        &domain.StateWaitingForLink,
	}

	if err := h.customerService.Update(ctx, customer.CustomerID, updateDTO); err != nil {
		return fmt.Errorf("customerRepo.Update: %w", err)
	}

	if customer.HasPromocode() {
		discount := customer.MustGetPromocode().GetClothingDiscount()
		if err := h.sendMessage(chatID, fmt.Sprintf("Стоимость товара с учетом скидки: %d ₽", priceRub-uint64(discount))); err != nil {
			return err
		}
	} else {
		if err := h.sendMessage(chatID, fmt.Sprintf("Стоимость товара: %d ₽", priceRub)); err != nil {
			return err
		}
	}

	if err := h.sendMessage(chatID, deliveryOnlyToMoscowTemplate); err != nil {
		return err
	}

	return h.sendMessage(chatID, askForLinkTemplate)
}

func (h *handler) HandleLinkInput(ctx context.Context, m *tg.Message) error {
	var (
		chatID = m.From.ID
		link   = strings.TrimSpace(m.Text)
	)

	// validate state
	customer, err := h.checkRequiredState(ctx, chatID, domain.StateWaitingForLink)
	if err != nil {
		return fmt.Errorf("check required state: %w", err)
	}

	if ok := url.IsValidDW4URL(link); !ok {
		if err := h.sendMessage(chatID, "Неправильная ссылка! Смотри инструкцию"); err != nil {
			return err
		}
		return h.sendMessage(chatID, "Введи повторно корректную ссылку 😀")
	}

	customer.UpdateLastEditPositionLink(link)
	customer.Cart.Add(*customer.LastEditPosition)
	updateDTO := dto.UpdateClothingCustomerDTO{
		LastPosition: customer.LastEditPosition,
		Cart:         &customer.Cart,
		State:        &domain.StateDefault,
	}

	if err := h.customerService.Update(ctx, customer.CustomerID, updateDTO); err != nil {
		return fmt.Errorf("customerRepo.Update: %w", err)
	}

	if err := h.sendMessage(chatID, fmt.Sprintf("Товар по ссылке: %s", link)); err != nil {
		return err
	}

	positionAddedMsg := tg.NewMessage(chatID, "Позиция успешно добавлена!")
	positionAddedMsg.ReplyMarkup = bottomMenuButtons
	return h.cleanSend(positionAddedMsg)
}
