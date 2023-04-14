package handler

import (
	"context"
	"fmt"
	"strings"

	"domain"
	"dto"
	"household_bot/internal/telegram/buttons"
	"household_bot/internal/telegram/callback"
	"household_bot/internal/telegram/templates"
	"household_bot/internal/telegram/tg_errors"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *handler) AskForOrderType(ctx context.Context, chatID int64) error {
	var telegramID = chatID

	if err := h.customerRepo.UpdateState(ctx, telegramID, domain.StateWaitingForOrderType); err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "AskForOrderType",
			CausedBy:    "UpdateState",
		})
	}
	return h.sendWithKeyboard(chatID, templates.AskForOrderType(), buttons.OrderTypeSelect)
}

func (h *handler) HandleOrderTypeInput(ctx context.Context, chatID int64, args []string) error {
	var (
		telegramID   = chatID
		orderTypeStr = args[0]
	)
	if err := h.checkRequiredState(ctx, chatID, domain.StateWaitingForOrderType); err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "HandleOrderTypeInput",
			CausedBy:    "checkRequiredState",
		})
	}

	orderType, err := domain.NewOrderTypeFromString(orderTypeStr)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "HandleOrderTypeInput",
			CausedBy:    "NewOrderTypeFromString",
		})
	}

	customer, err := h.customerRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "HandleOrderTypeInput",
			CausedBy:    "GetByTelegramID",
		})
	}

	customer.UpdateMetaOrderType(orderType)
	updateDTO := dto.UpdateHouseholdCustomerDTO{
		State: &domain.StateWaitingForFIO,
		Meta:  &customer.Meta,
	}
	if err := h.customerRepo.Update(ctx, customer.CustomerID, updateDTO); err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "HandleOrderTypeInput",
			CausedBy:    "Update",
		})
	}
	if err := h.askForFIO(ctx, chatID); err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "HandleOrderTypeInput",
			CausedBy:    "askForFIO",
		})
	}
	return nil
}

func (h *handler) askForFIO(ctx context.Context, chatID int64) error {
	var telegramID = chatID
	if err := h.customerRepo.UpdateState(ctx, telegramID, domain.StateWaitingForFIO); err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "askForFIO",
			CausedBy:    "UpdateState",
		})
	}
	return h.sendMessage(chatID, templates.AskForFIO())
}

func (h *handler) HandleFIOInput(ctx context.Context, m *tg.Message) error {
	var (
		chatID     = m.From.ID
		telegramID = chatID
		fullName   = strings.TrimSpace(m.Text)
	)

	if err := h.checkRequiredState(ctx, chatID, domain.StateWaitingForFIO); err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "HandleFIOInput",
			CausedBy:    "checkRequiredState",
		})
	}

	customer, err := h.customerRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "HandleFIOInput",
			CausedBy:    "GetByTelegramID",
		})
	}

	if !domain.IsValidFullName(fullName) {
		if err := h.sendMessage(chatID, templates.InvalidFIO()); err != nil {
			return tg_errors.New(tg_errors.Config{
				OriginalErr: err,
				Handler:     "HandleFIOInput",
				CausedBy:    "sendMessage",
			})
		}
	}

	updateDTO := dto.UpdateHouseholdCustomerDTO{
		State:    &domain.StateWaitingForPhoneNumber,
		FullName: &fullName,
	}
	if err := h.customerRepo.Update(ctx, customer.CustomerID, updateDTO); err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "HandleFIOInput",
			CausedBy:    "Update",
		})
	}

	if err := h.sendMessage(chatID, fmt.Sprintf("Спасибо, %s. ", fullName)); err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "HandleFIOInput",
			CausedBy:    "sendMessage",
		})
	}

	return h.sendMessage(chatID, templates.AskForPhoneNumber())
}

func (h *handler) HandlePhoneNumberInput(ctx context.Context, m *tg.Message) error {
	var (
		chatID      = m.From.ID
		telegramID  = chatID
		phoneNumber = strings.TrimSpace(m.Text)
	)

	if err := h.checkRequiredState(ctx, chatID, domain.StateWaitingForPhoneNumber); err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "HandlePhoneNumberInput",
			CausedBy:    "checkRequiredState",
		})
	}

	if !domain.IsValidPhoneNumber(phoneNumber) {
		if err := h.sendMessage(chatID, "Неправильный формат номера телефона.\n"+templates.AskForPhoneNumber()); err != nil {
			return tg_errors.New(tg_errors.Config{
				OriginalErr: err,
				Handler:     "HandlePhoneNumberInput",
				CausedBy:    "sendMessage",
			})
		}
	}

	customer, err := h.customerRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "HandlePhoneNumberInput",
			CausedBy:    "GetByTelegramID",
		})
	}

	updateDTO := dto.UpdateHouseholdCustomerDTO{
		State:       &domain.StateWaitingForDeliveryAddress,
		PhoneNumber: &phoneNumber,
	}
	if err := h.customerRepo.Update(ctx, customer.CustomerID, updateDTO); err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "HandlePhoneNumberInput",
			CausedBy:    "Update",
		})
	}

	if err := h.sendMessage(chatID, fmt.Sprintf("Спасибо, номер [%s] принят!", phoneNumber)); err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "HandlePhoneNumberInput",
			CausedBy:    "sendMessage",
		})
	}

	return h.sendMessage(chatID, templates.AskForDeliveryAddress())
}

func (h *handler) HandleDeliveryAddressInput(ctx context.Context, m *tg.Message) error {
	var (
		chatID     = m.From.ID
		telegramID = chatID
		address    = strings.TrimSpace(m.Text)
	)

	if err := h.checkRequiredState(ctx, chatID, domain.StateWaitingForDeliveryAddress); err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "HandleDeliveryAddressInput",
			CausedBy:    "checkRequiredState",
		})
	}

	customer, err := h.customerRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "HandleDeliveryAddressInput",
			CausedBy:    "GetByTelegramID",
		})
	}

	shortID, err := h.orderService.GetFreeShortID(ctx)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "HandleDeliveryAddressInput",
			CausedBy:    "GetFreeShortID",
		})
	}

	isExpress := *customer.Meta.NextOrderType == domain.OrderTypeExpress
	order := domain.NewHouseholdOrder(customer, address, isExpress, shortID)

	if err := h.orderService.Save(ctx, order); err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "HandleDeliveryAddressInput",
			CausedBy:    "Save",
		})
	}

	customer.Cart.Clear()
	updateDTO := dto.UpdateHouseholdCustomerDTO{
		Cart:  &customer.Cart,
		State: &domain.StateDefault,
	}
	if err := h.customerRepo.Update(ctx, customer.CustomerID, updateDTO); err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "HandleDeliveryAddressInput",
			CausedBy:    "Update",
		})
	}

	if err := h.sendOrder(ctx, order, chatID); err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "HandleDeliveryAddressInput",
			CausedBy:    "sendOrderPreview",
		})
	}
	return nil
}

func (h *handler) sendOrder(ctx context.Context,
	order domain.HouseholdOrder,
	chatID int64) error {

	orderText := templates.RenderOrder(order)
	if err := h.sendMessage(chatID, orderText); err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "sendOrder",
			CausedBy:    "sendMessage",
		})
	}

	updateDTO := dto.UpdateHouseholdCustomerDTO{
		Meta:  &domain.Meta{},
		State: &domain.StateDefault,
	}
	if err := h.customerRepo.Update(ctx, order.Customer.CustomerID, updateDTO); err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "sendOrder",
			CausedBy:    "Update",
		})
	}

	requisitesMsg := tg.NewMessage(chatID, templates.Requisites(order.ShortID, domain.AdminRequisites))
	sentRequisitesMsg, err := h.bot.Send(requisitesMsg)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "Send",
			CausedBy:    "Update",
		})
	}

	editButton := tg.NewEditMessageReplyMarkup(chatID,
		sentRequisitesMsg.MessageID,
		buttons.NewPaymentButton(callback.AcceptPayment, order.ShortID),
	)

	return h.cleanSend(editButton)
}

func (h *handler) HandlePayment(ctx context.Context, c *tg.CallbackQuery, args []string) error {
	var (
		chatID       = c.From.ID
		telegramID   = chatID
		shortOrderID = args[0]
	)

	customer, err := h.customerRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "HandlePayment",
			CausedBy:    "GetByTelegramID",
		})
	}

	if err := h.orderService.UpdateToPaid(ctx, customer.CustomerID, shortOrderID); err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "HandlePayment",
			CausedBy:    "UpdateToPaid",
		})
	}

	editButtons := tg.NewEditMessageReplyMarkup(
		chatID,
		c.Message.MessageID,
		buttons.NewSuccessfulPaymentButton(shortOrderID),
	)
	if err := h.cleanSend(editButtons); err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "HandlePayment",
			CausedBy:    "cleanSend",
		})
	}

	return h.sendWithKeyboard(
		chatID,
		templates.SuccessfulPayment(*customer.FullName, shortOrderID),
		buttons.MakeOrder,
	)
}
