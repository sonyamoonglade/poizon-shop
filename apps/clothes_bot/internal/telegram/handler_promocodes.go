package telegram

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"domain"
	"dto"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *handler) AskForPromocode(ctx context.Context, chatID int64) error {
	customer, err := h.customerService.GetByTelegramID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get by telegram id: %w", err)
	}
	if customer.HasPromocode() {
		return h.sendMessage(chatID, "Ты уже вводил промокод!")
	}
	if err := h.customerService.UpdateState(ctx, chatID, domain.StateWaitingForPromocode); err != nil {
		return fmt.Errorf("update state: %w", err)
	}
	if err := h.sendMessage(chatID, promocodeWarning()); err != nil {
		return err
	}
	return h.sendMessage(chatID, askForPromocode())
}

func (h *handler) HandlePromocodeInput(ctx context.Context, m *tg.Message) error {
	var (
		chatID       = m.Chat.ID
		promoShortID = strings.TrimSpace(m.Text) // shortID of domain.Promocode
	)
	customer, err := h.checkRequiredState(ctx, chatID, domain.StateWaitingForPromocode)
	if err != nil {
		return fmt.Errorf("check required state: %w", err)
	}

	promo, err := h.promocodeRepo.GetByShortID(ctx, promoShortID)
	if err != nil {
		if errors.Is(err, domain.ErrNoPromocode) {
			return h.sendMessage(chatID, fmt.Sprintf("Промокод %s не существует!", promoShortID))
		}
		return fmt.Errorf("promo get by short id: %w", err)
	}

	if customer.HasPromocode() {
		return h.sendMessage(chatID, "Ты уже вводил промокод!")
	}

	customer.UsePromocode(promo)

	err = h.customerService.Update(ctx, customer.CustomerID, dto.UpdateClothingCustomerDTO{
		PromocodeID: customer.PromocodeID,
	})
	if err != nil {
		return fmt.Errorf("customer service update: %w", err)
	}

	return h.sendMessage(chatID, promocodeUseSuccess(promo.ShortID, promo.GetClothingDiscount()))
}
