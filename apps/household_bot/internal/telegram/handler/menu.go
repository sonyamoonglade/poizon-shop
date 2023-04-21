package handler

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"dto"
	"household_bot/internal/telegram/callback"

	"domain"
	"household_bot/internal/telegram/buttons"
	"household_bot/internal/telegram/templates"
	"household_bot/internal/telegram/tg_errors"
	"household_bot/pkg/telegram"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *handler) Start(ctx context.Context, m *tg.Message) error {
	var (
		telegramID = m.Chat.ID
		username   = domain.MakeUsername(m.From.String())
	)
	_, err := h.customerRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		if !errors.Is(err, domain.ErrCustomerNotFound) {
			return tg_errors.New(tg_errors.Config{
				OriginalErr: err,
				Handler:     "Start",
				CausedBy:    "GetByTelegramID",
			})

		}

		err := h.customerRepo.Save(ctx, domain.NewHouseholdCustomer(telegramID, username))
		if err != nil {
			return tg_errors.New(tg_errors.Config{
				OriginalErr: err,
				Handler:     "Start",
				CausedBy:    "Save",
			})
		}
	}
	return h.sendWithKeyboard(telegramID, templates.StartGreeting(username), buttons.Start)
}

func (h *handler) Menu(ctx context.Context, chatID int64, deleteMsgID *int) error {
	if deleteMsgID != nil {
		if err := h.bot.CleanRequest(tg.NewDeleteMessage(chatID, *deleteMsgID)); err != nil {
			return err
		}
	}
	customer, err := h.customerRepo.GetByTelegramID(ctx, chatID)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "Menu",
			CausedBy:    "GetByTelegramID",
		})
	}
	if err := h.customerRepo.UpdateState(ctx, chatID, domain.StateDefault); err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "Menu",
			CausedBy:    "UpdateState",
		})
	}
	showPromoButton := !customer.HasPromocode()
	return h.sendWithKeyboard(chatID, "Меню", buttons.Menu(showPromoButton))
}

func (h *handler) Catalog(ctx context.Context, chatID int64, prevMsgID *int) error {
	var c tg.Chattable
	if prevMsgID != nil {
		editMsg := tg.NewEditMessageText(chatID, *prevMsgID, "Выбери тип товара")
		editMsg.ReplyMarkup = &buttons.CatalogType
		c = editMsg
	} else {
		msg := tg.NewMessage(chatID, "Выбери тип товара")
		msg.ReplyMarkup = buttons.CatalogType
		c = msg
	}

	return h.sendWithMessageID(c, func(msgID int) error {
		catalogMsg := telegram.CatalogMsg{
			MsgID:  msgID,
			ChatID: chatID,
		}
		err := h.catalogMsgService.Save(ctx, catalogMsg)
		if err != nil {
			if errors.Is(err, telegram.ErrMessageAlreadyExists) {
				return nil
			}
			return tg_errors.New(tg_errors.Config{
				OriginalErr: err,
				Handler:     "Catalog",
				CausedBy:    "Save",
			})
		}

		return nil
	})
}

func (h *handler) MyOrders(ctx context.Context, chatID int64) error {
	var telegramID = chatID

	customer, err := h.customerRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return err
	}

	orders, err := h.orderService.GetAllForCustomer(ctx, customer.CustomerID, domain.SourceHousehold)
	if err != nil {
		if errors.Is(err, domain.ErrNoOrders) {
			return h.sendMessage(chatID, "У тебя еще нет заказов 🦕")
		}
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "MyOrders",
			CausedBy:    "GetAllForCustomer",
		})
	}
	if orders == nil {
		return h.sendMessage(chatID, "У тебя еще нет заказов 🦕")
	}

	var name string
	if customer.FullName != nil {
		name = *customer.FullName
	} else {
		name = *customer.Username
	}

	return h.sendMessage(chatID, templates.RenderMyOrders(name, orders))
}

func (h *handler) AskForISBN(ctx context.Context, chatID int64) error {
	if err := h.customerRepo.UpdateState(ctx, chatID, domain.StateWaitingForISBN); err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "AskForISBN",
			CausedBy:    "UpdateState",
		})
	}
	return h.sendMessage(chatID, "Введи артикул товара: ")
}

func (h *handler) HandleProductByISBN(ctx context.Context, m *tg.Message) error {
	var (
		chatID = m.Chat.ID
		isbn   = strings.TrimSpace(m.Text)
	)

	product, ok := h.catalogProvider.GetProductByISBN(isbn)
	if !ok {
		if err := h.sendMessage(chatID, fmt.Sprintf("Товар с артикулом: %s не существует", isbn)); err != nil {
			return err
		}
		if err := h.AskForISBN(ctx, chatID); err != nil {
			return tg_errors.New(tg_errors.Config{
				OriginalErr: err,
				Handler:     "HandleProductByISBN",
				CausedBy:    "AskForISBN",
			})
		}
		return nil
	}

	customer, err := h.customerRepo.GetByTelegramID(ctx, chatID)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "HandleProductByISBN",
			CausedBy:    "GetByTelegramID",
		})
	}

	backButton := buttons.NewBackButton(callback.Menu, nil, nil, nil)
	err = h.renderProductCard(ctx, chatID, product, customer, buttons.NewISBNProductCardButtons(isbn, backButton))
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "HandleProductByISBN",
			CausedBy:    "renderProductCard",
		})
	}
	return nil
}

func (h *handler) AskForPromocode(ctx context.Context, chatID int64) error {
	if err := h.customerRepo.UpdateState(ctx, chatID, domain.StateWaitingForPromocode); err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "AskForPromocode",
			CausedBy:    "UpdateState",
		})
	}
	if err := h.sendMessage(chatID, templates.PromocodeWarning()); err != nil {
		return err
	}
	return h.sendMessage(chatID, templates.AskForPromocode())
}

func (h *handler) HandlePromocodeInput(ctx context.Context, m *tg.Message) error {
	var (
		chatID       = m.Chat.ID
		promoShortID = strings.TrimSpace(m.Text) // shortID of domain.Promocode
	)
	customer, err := h.checkRequiredState(ctx, chatID, domain.StateWaitingForPromocode)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "HandlePromocodeInput",
			CausedBy:    "checkRequiredState",
		})
	}

	promo, err := h.promocodeRepo.GetByShortID(ctx, promoShortID)
	if err != nil {
		if errors.Is(err, domain.ErrNoPromocode) {
			return h.sendMessage(chatID, fmt.Sprintf("Промокод %s не существует!", promoShortID))
		}
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "HandlePromocodeInput",
			CausedBy:    "GetByShortID",
		})
	}

	if customer.HasPromocode() {
		return h.sendMessage(chatID, "Вы уже вводили промокод!")
	}

	customer.UsePromocode(promo)

	err = h.customerRepo.Update(ctx, customer.CustomerID, dto.UpdateHouseholdCustomerDTO{
		PromocodeID: customer.PromocodeID,
	})
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "HandlePromocodeInput",
			CausedBy:    "Update",
		})
	}

	return h.sendMessage(chatID, templates.PromocodeUseSuccess(promo.ShortID, promo.GetDiscount(domain.SourceHousehold)))
}
