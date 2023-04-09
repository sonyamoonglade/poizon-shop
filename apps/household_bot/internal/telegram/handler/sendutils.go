package handler

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *handler) sendWithKeyboard(chatID int64, text string, keyboard interface{}) error {
	m := tg.NewMessage(chatID, text)
	m.ReplyMarkup = keyboard
	return h.cleanSend(m)
}

func (h *handler) cleanSend(c tg.Chattable) error {
	_, err := h.bot.Send(c)
	return err
}

//func (h *handler) checkRequiredState(ctx context.Context, want domain.State, telegramID int64) error {
//	// customer, err := h.customerRepo.GetByTelegramID(ctx, telegramID)
//	// if err != nil {
//	// 	return fmt.Errorf("checkRequiredState: %w", err)
//	// }
//	// if customer.TgState != want {
//	// 	return ErrInvalidState
//	// }
//	// return nil
//}

func (h *handler) sendMessage(chatID int64, text string) error {
	return h.cleanSend(tg.NewMessage(chatID, text))
}
