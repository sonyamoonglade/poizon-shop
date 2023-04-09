package telegram

import (
	"context"
	"fmt"
	"strings"

	"domain"
	"dto"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *handler) AskForFIO(ctx context.Context, chatID int64) error {
	var telegramID = chatID

	if err := h.customerRepo.UpdateState(ctx, telegramID, domain.StateWaitingForFIO); err != nil {
		return err
	}
	return h.sendMessage(chatID, askForFIOTemplate)
}

func (h *handler) HandleFIOInput(ctx context.Context, m *tg.Message) error {
	var (
		chatID     = m.From.ID
		telegramID = chatID
		fullName   = strings.TrimSpace(m.Text)
	)

	if err := h.checkRequiredState(ctx, domain.StateWaitingForFIO, chatID); err != nil {
		return err
	}

	customer, err := h.customerRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return fmt.Errorf("customerRepo.GetByTelegramID: %w", err)
	}

	if !domain.IsValidFullName(fullName) {
		return h.sendMessage(chatID, invalidFIOInputTemplate)
	}

	updateDTO := dto.UpdateClothingCustomerDTO{
		State:    &domain.StateWaitingForPhoneNumber,
		FullName: &fullName,
	}

	if err := h.customerRepo.Update(ctx, customer.CustomerID, updateDTO); err != nil {
		return fmt.Errorf("customerRepo.Update: %w", err)
	}

	if err := h.sendMessage(chatID, fmt.Sprintf("Спасибо, %s. ", fullName)); err != nil {
		return err
	}

	return h.sendMessage(chatID, askForPhoneNumberTemplate)
}

func (h *handler) HandlePhoneNumberInput(ctx context.Context, m *tg.Message) error {
	var (
		chatID      = m.From.ID
		telegramID  = chatID
		phoneNumber = strings.TrimSpace(m.Text)
	)

	if err := h.checkRequiredState(ctx, domain.StateWaitingForPhoneNumber, chatID); err != nil {
		return err
	}

	if !domain.IsValidPhoneNumber(phoneNumber) {
		return h.sendMessage(chatID, "Неправильный формат номера телефона.\n"+askForPhoneNumberTemplate)
	}

	customer, err := h.customerRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return fmt.Errorf("customerRepo.GetByTelegramID: %w", err)
	}

	updateDTO := dto.UpdateClothingCustomerDTO{
		State:       &domain.StateWaitingForDeliveryAddress,
		PhoneNumber: &phoneNumber,
	}

	if err := h.customerRepo.Update(ctx, customer.CustomerID, updateDTO); err != nil {
		return fmt.Errorf("customerRepo.Update: %w", err)
	}

	if err := h.sendMessage(chatID, fmt.Sprintf("Спасибо, номер [%s] принят!", phoneNumber)); err != nil {
		return err
	}

	return h.sendMessage(chatID, askForDeliveryAddressTemplate)
}

func (h *handler) HandleDeliveryAddressInput(ctx context.Context, m *tg.Message) error {
	var (
		chatID     = m.From.ID
		telegramID = chatID
		address    = strings.TrimSpace(m.Text)
	)

	// validate state
	if err := h.checkRequiredState(ctx, domain.StateWaitingForDeliveryAddress, chatID); err != nil {
		return err
	}

	customer, err := h.customerRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return fmt.Errorf("customerRepo.GetByTelegramID: %w", err)
	}

	updateDTO := dto.UpdateClothingCustomerDTO{
		LastPosition: &domain.ClothingPosition{},
		Cart:         &domain.ClothingCart{},
		State:        &domain.StateDefault,
	}

	if err := h.customerRepo.Update(ctx, customer.CustomerID, updateDTO); err != nil {
		return fmt.Errorf("customerRepo.Update: %w", err)
	}

	shortID, err := h.orderService.GetFreeShortID(ctx)
	if err != nil {
		return err
	}

	isExpress := *customer.Meta.NextOrderType == domain.OrderTypeExpress
	order := domain.NewOrder(customer, address, isExpress, shortID, domain.SourceClothing)

	if err := h.orderService.Save(ctx, order); err != nil {
		return err
	}

	return h.prepareOrderPreview(ctx, customer, order, chatID)
}

func (h *handler) prepareOrderPreview(ctx context.Context, customer domain.Customer, order domain.ClothingOrder, chatID int64) error {
	out := getOrderStart(orderStartArgs{
		fullName:        *customer.FullName,
		shortOrderID:    order.ShortID,
		phoneNumber:     *customer.PhoneNumber,
		isExpress:       order.IsExpress,
		deliveryAddress: order.DeliveryAddress,
		nCartItems:      len(order.Cart),
	})

	for i, cartItem := range order.Cart {
		out += getPositionTemplate(cartPositionPreviewArgs{
			n:         i + 1,
			link:      cartItem.ShopLink,
			size:      cartItem.Size,
			category:  string(cartItem.Category),
			priceRub:  cartItem.PriceRUB,
			priceYuan: cartItem.PriceYUAN,
		})
	}

	out += getOrderEnd(order.AmountRUB)

	if err := h.sendMessage(chatID, out); err != nil {
		return err
	}

	updateDTO := dto.UpdateClothingCustomerDTO{
		Meta:  &domain.Meta{},
		Cart:  new(domain.ClothingCart),
		State: &domain.StateDefault,
	}

	if err := h.customerRepo.Update(ctx, customer.CustomerID, updateDTO); err != nil {
		return err
	}

	requisitesMsg := tg.NewMessage(chatID, getRequisites(domain.AdminRequisites, order.ShortID))
	sentRequisitesMsg, err := h.b.Send(requisitesMsg)
	if err != nil {
		return err
	}

	editButton := tg.NewEditMessageReplyMarkup(chatID, sentRequisitesMsg.MessageID, preparePaymentButton(order.ShortID))
	return h.cleanSend(editButton)
}

func (h *handler) HandlePayment(ctx context.Context, shortOrderID string, c *tg.CallbackQuery) error {
	var (
		chatID     = c.From.ID
		telegramID = chatID
	)

	customer, err := h.customerRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return fmt.Errorf("customerRepo.GetByTelegramID: %w", err)
	}

	if err := h.orderService.UpdateToPaid(ctx, customer.CustomerID, shortOrderID); err != nil {
		return err
	}

	editButtons := tg.NewEditMessageReplyMarkup(chatID, c.Message.MessageID, prepareAfterPaidButtons(shortOrderID))
	if err := h.cleanSend(editButtons); err != nil {
		return err
	}

	return h.sendWithKeyboard(chatID, getAfterPaid(*customer.FullName, shortOrderID), makeOrderButtons)
}
