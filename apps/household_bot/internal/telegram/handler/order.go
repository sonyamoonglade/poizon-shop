package handler

import (
	"context"
	"fmt"
	"strings"

	"domain"
	"dto"
	"household_bot/internal/telegram/buttons"
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

	if err := h.sendOrderPreview(ctx, customer, order, chatID); err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "HandleDeliveryAddressInput",
			CausedBy:    "sendOrderPreview",
		})
	}
	return nil
}

func (h *handler) sendOrderPreview(ctx context.Context,
	customer domain.HouseholdCustomer,
	order domain.HouseholdOrder,
	chatID int64) error {
	// out := getOrderStart(orderStartArgs{
	// 	fullName:        *customer.FullName,
	// 	shortOrderID:    order.ShortID,
	// 	phoneNumber:     *customer.PhoneNumber,
	// 	isExpress:       order.IsExpress,
	// 	deliveryAddress: order.DeliveryAddress,
	// 	nCartItems:      len(order.Cart),
	// })

	// for i, cartItem := range order.Cart {
	// 	out += getPositionTemplate(cartPositionPreviewArgs{
	// 		n:         i + 1,
	// 		link:      cartItem.ShopLink,
	// 		size:      cartItem.Size,
	// 		category:  string(cartItem.Category),
	// 		priceRub:  cartItem.PriceRUB,
	// 		priceYuan: cartItem.PriceYUAN,
	// 	})
	// }

	// out += getOrderEnd(order.AmountRUB)

	// if err := h.sendMessage(chatID, out); err != nil {
	// 	return err
	// }
	// updateDTO := dto.UpdateHouseholdCustomerDTO{
	// 	Meta:  &domain.Meta{},
	// 	State: &domain.StateDefault,
	// }

	// if err := h.customerRepo.Update(ctx, customer.CustomerID, updateDTO); err != nil {
	// 	return err
	// }

	// requisitesMsg := tg.NewMessage(chatID, templates.Requisites(order.ShortID, domain.AdminRequisites))
	// sentRequisitesMsg, err := h.bot.Send(requisitesMsg)
	// if err != nil {
	// 	return err
	// }

	// editButton := tg.NewEditMessageReplyMarkup(chatID,
	// 	sentRequisitesMsg.MessageID,
	// 	buttons.NewPaymentButton(callback.AcceptPayment, order.ShortID))
	// return h.cleanSend(editButton)
	return nil
}

func (h *handler) HandlePayment(ctx context.Context, shortOrderID string, c *tg.CallbackQuery) error {
	return nil
	// var (
	// 	chatID     = c.From.ID
	// 	telegramID = chatID
	// )

	// customer, err := h.customerRepo.GetByTelegramID(ctx, telegramID)
	// if err != nil {
	// 	return fmt.Errorf("customerRepo.GetByTelegramID: %w", err)
	// }

	// if err := h.orderService.UpdateToPaid(ctx, customer.CustomerID, shortOrderID); err != nil {
	// 	return err
	// }

	// editButtons := tg.NewEditMessageReplyMarkup(chatID, c.Message.MessageID, prepareAfterPaidButtons(shortOrderID))
	// if err := h.cleanSend(editButtons); err != nil {
	// 	return err
	// }

	// return h.sendWithKeyboard(chatID, getAfterPaid(*customer.FullName, shortOrderID), makeOrderButtons)
}
