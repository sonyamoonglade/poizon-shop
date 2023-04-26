package handler

import (
	"context"
	"fmt"
	"strings"

	"domain"
	"dto"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"household_bot/internal/telegram/buttons"
	"household_bot/internal/telegram/callback"
	"household_bot/internal/telegram/templates"
	"household_bot/internal/telegram/tg_errors"
	"utils/ptr"
)

func (h *handler) AskForFIO(ctx context.Context, chatID int64) error {
	telegramID := chatID
	// Beforehand check the cart for validity
	if err := h.sendMessage(chatID, templates.CheckingCart()); err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "AskForFIO",
			CausedBy:    "sendMessage",
		})
	}

	customer, err := h.customerService.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "AskForFIO",
			CausedBy:    "GetByTelegramID",
		})
	}

	// Perform a *long* check that checks whether all products in cart exist
	first, _ := customer.Cart.First()
	someCategoryID := first.CategoryID

	category, err := h.categoryService.GetByID(ctx, someCategoryID)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "AskForFIO",
			CausedBy:    "GetByID",
		})
	}

	ok, missingProduct, err := h.categoryService.CheckIfAllProductsExist(ctx, customer.Cart, category.InStock)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "AskForFIO",
			CausedBy:    "CheckIfAllProductsExist",
		})
	}

	if !ok {
		err := h.sendWithKeyboard(
			chatID,
			templates.ProductNotFound(
				missingProduct.Name,
				missingProduct.ISBN,
			),
			buttons.GotoCart,
		)

		if err != nil {
			return tg_errors.New(tg_errors.Config{
				OriginalErr: err,
				Handler:     "AskForFIO",
				CausedBy:    "sendMessage",
			})
		}
		return nil
	}

	if err := h.sendMessage(chatID, "Все хорошо, можешь создавать заказ"); err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "AskForFIO",
			CausedBy:    "sendMessage",
		})
	}

	if err := h.customerService.UpdateState(ctx, telegramID, domain.StateWaitingForFIO); err != nil {
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
		chatID   = m.From.ID
		fullName = strings.TrimSpace(m.Text)
	)
	customer, err := h.checkRequiredState(ctx, chatID, domain.StateWaitingForFIO)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "HandleFIOInput",
			CausedBy:    "checkRequiredState",
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
		return nil
	}

	updateDTO := dto.UpdateHouseholdCustomerDTO{
		State:    &domain.StateWaitingForPhoneNumber,
		FullName: &fullName,
	}
	if err := h.customerService.Update(ctx, customer.CustomerID, updateDTO); err != nil {
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
		phoneNumber = strings.TrimSpace(m.Text)
	)

	customer, err := h.checkRequiredState(ctx, chatID, domain.StateWaitingForPhoneNumber)
	if err != nil {
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
		return nil
	}

	updateDTO := dto.UpdateHouseholdCustomerDTO{
		State:       &domain.StateWaitingForDeliveryAddress,
		PhoneNumber: &phoneNumber,
	}
	if err := h.customerService.Update(ctx, customer.CustomerID, updateDTO); err != nil {
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
		chatID  = m.From.ID
		address = strings.TrimSpace(m.Text)
	)

	customer, err := h.checkRequiredState(ctx, chatID, domain.StateWaitingForDeliveryAddress)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "HandleDeliveryAddressInput",
			CausedBy:    "checkRequiredState",
		})
	}

	firstProduct, ok := customer.Cart.First()
	if !ok {
		return h.sendMessage(chatID, "Твоя корзина пуста!")
	}
	category, ok := h.catalogProvider.GetCategoryByID(firstProduct.CategoryID)
	if !ok {
		return h.handleIfCategoryNotFound(ctx, chatID, customer, firstProduct)
	}

	order, err := h.makeOrderUsecase.NewOrder(ctx, address, customer, category.InStock)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "HandleDeliveryAddressInput",
			CausedBy:    "NewOrder",
		})
	}

	err = h.sendOrder(
		ctx,
		chatID,
		order,
		customer.HasPromocode(),
		ptr.Ptr(customer.MustGetPromocode().GetHouseholdDiscount()),
	)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "HandleDeliveryAddressInput",
			CausedBy:    "sendOrder",
		})
	}
	return nil
}

func (h *handler) sendOrder(ctx context.Context, chatID int64, order domain.HouseholdOrder, discounted bool, discount *uint32) error {
	if discounted && discount != nil {
		if err := h.sendMessage(chatID, templates.RenderOrderAfterPaymentWithDiscount(order, *discount)); err != nil {
			return err
		}
	} else {
		if err := h.sendMessage(chatID, templates.RenderOrderAfterPayment(order)); err != nil {
			return err
		}
	}

	updateDTO := dto.UpdateHouseholdCustomerDTO{
		State: &domain.StateDefault,
	}
	if err := h.customerService.Update(ctx, order.Customer.CustomerID, updateDTO); err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "sendOrder",
			CausedBy:    "Update",
		})
	}

	requisites := tg.NewMessage(chatID, templates.Requisites(order.ShortID, domain.AdminRequisites))
	return h.sendWithMessageID(requisites, func(msgID int) error {
		editButton := tg.NewEditMessageReplyMarkup(
			chatID,
			msgID,
			buttons.NewPaymentButton(callback.AcceptPayment, order.ShortID),
		)
		return h.cleanSend(editButton)
	})
}

func (h *handler) HandlePayment(ctx context.Context, c *tg.CallbackQuery, args []string) error {
	var (
		chatID       = c.From.ID
		telegramID   = chatID
		shortOrderID = args[0]
	)

	customer, err := h.customerService.GetByTelegramID(ctx, telegramID)
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
