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

	if err := h.customerService.UpdateState(ctx, telegramID, domain.StateWaitingForFIO); err != nil {
		return err
	}
	return h.sendMessage(chatID, askForFIOTemplate)
}

func (h *handler) HandleFIOInput(ctx context.Context, m *tg.Message) error {
	var (
		chatID   = m.From.ID
		fullName = strings.TrimSpace(m.Text)
	)

	customer, err := h.checkRequiredState(ctx, chatID, domain.StateWaitingForFIO)
	if err != nil {
		return fmt.Errorf("check required state: %w", err)
	}

	if !domain.IsValidFullName(fullName) {
		return h.sendMessage(chatID, invalidFIOInputTemplate)
	}

	updateDTO := dto.UpdateClothingCustomerDTO{
		State:    &domain.StateWaitingForPhoneNumber,
		FullName: &fullName,
	}

	if err := h.customerService.Update(ctx, customer.CustomerID, updateDTO); err != nil {
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
		phoneNumber = strings.TrimSpace(m.Text)
	)

	customer, err := h.checkRequiredState(ctx, chatID, domain.StateWaitingForPhoneNumber)
	if err != nil {
		return fmt.Errorf("check required state: %w", err)
	}

	if !domain.IsValidPhoneNumber(phoneNumber) {
		return h.sendMessage(chatID, "Неправильный формат номера телефона.\n"+askForPhoneNumberTemplate)
	}

	updateDTO := dto.UpdateClothingCustomerDTO{
		State:       &domain.StateWaitingForDeliveryAddress,
		PhoneNumber: &phoneNumber,
	}

	if err := h.customerService.Update(ctx, customer.CustomerID, updateDTO); err != nil {
		return fmt.Errorf("customerRepo.Update: %w", err)
	}

	if err := h.sendMessage(chatID, fmt.Sprintf("Спасибо, номер [%s] принят!", phoneNumber)); err != nil {
		return err
	}

	return h.sendMessage(chatID, askForDeliveryAddressTemplate)
}

func (h *handler) HandleDeliveryAddressInput(ctx context.Context, m *tg.Message) error {
	var (
		chatID  = m.From.ID
		address = strings.TrimSpace(m.Text)
	)

	// validate state
	customer, err := h.checkRequiredState(ctx, chatID, domain.StateWaitingForDeliveryAddress)
	if err != nil {
		return fmt.Errorf("check required state: %w", err)
	}

	shortID, err := h.orderService.GetFreeShortID(ctx)
	if err != nil {
		return err
	}

	isExpress := customer.Meta.NextOrderType.IsExpress()
	order := domain.NewClothingOrder(customer, address, isExpress, shortID)

	if err := h.orderService.Save(ctx, order); err != nil {
		return err
	}

	updateDTO := dto.UpdateClothingCustomerDTO{
		LastPosition: &domain.ClothingPosition{},
		Cart:         &domain.ClothingCart{},
		State:        &domain.StateDefault,
	}

	if err := h.customerService.Update(ctx, customer.CustomerID, updateDTO); err != nil {
		return fmt.Errorf("customerRepo.Update: %w", err)
	}

	return h.prepareOrderPreview(ctx, customer, order, chatID)
}

func (h *handler) prepareOrderPreview(ctx context.Context, customer domain.ClothingCustomer, order domain.ClothingOrder, chatID int64) error {
	out := getOrderStart(orderStartArgs{
		fullName:        *customer.FullName,
		shortOrderID:    order.ShortID,
		phoneNumber:     *customer.PhoneNumber,
		isExpress:       order.IsExpress,
		deliveryAddress: order.DeliveryAddress,
		nCartItems:      len(order.Cart),
	})
	var (
		discount   = new(uint32)
		discounted bool
	)
	if customer.HasPromocode() {
		discounted = true
		*discount = customer.MustGetPromocode().GetClothingDiscount()
	}
	for i, cartItem := range order.Cart {
		if discounted {
			out += getDiscountedPositionTemplate(cartPositionPreviewDiscountedArgs{
				n:           i + 1,
				link:        cartItem.ShopLink,
				size:        cartItem.Size,
				discountRub: *discount,
				category:    string(cartItem.Category),
				priceRub:    cartItem.PriceRUB,
				priceYuan:   cartItem.PriceYUAN,
			})
			continue
		}
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

	if err := h.customerService.Update(ctx, customer.CustomerID, updateDTO); err != nil {
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

	customer, err := h.customerService.GetByTelegramID(ctx, telegramID)
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
