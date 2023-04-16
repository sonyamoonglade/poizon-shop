package handler

import (
	"context"
	"fmt"

	"domain"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/multierr"
)

func (h *handler) sendWithKeyboard(chatID int64, text string, keyboard interface{}) error {
	m := tg.NewMessage(chatID, text)
	m.ReplyMarkup = keyboard
	return h.cleanSend(m)
}

func (h *handler) sendWithMessageID(c tg.Chattable, f func(msgID int) error) error {
	msg, err := h.bot.Send(c)
	if err != nil {
		return err
	}
	return f(msg.MessageID)
}

func (h *handler) cleanSend(c tg.Chattable) error {
	_, err := h.bot.Send(c)
	return err
}
func (h *handler) sendBulk(cs ...tg.Chattable) error {
	var errors error
	for _, c := range cs {
		if _, err := h.bot.Send(c); err != nil {
			errors = multierr.Append(errors, fmt.Errorf("sendBulk: %w", err))
		}
	}
	return errors
}

func (h *handler) checkRequiredState(ctx context.Context, telegramID int64, want domain.State) error {
	customer, err := h.customerRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return fmt.Errorf("checkRequiredState: %w", err)
	}
	if customer.State != want {
		return ErrInvalidState
	}
	return nil
}

func (h *handler) sendMessage(chatID int64, text string) error {
	return h.cleanSend(tg.NewMessage(chatID, text))
}
