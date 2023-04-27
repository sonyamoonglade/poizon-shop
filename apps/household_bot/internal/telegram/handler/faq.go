package handler

import (
	"context"
	"strconv"

	"household_bot/internal/telegram/buttons"
	"household_bot/internal/telegram/templates"
	"household_bot/internal/telegram/tg_errors"
)

func (h *handler) FAQ(ctx context.Context, chatID int64) error {
	return h.sendWithKeyboard(chatID, "❔Часто задаваемые вопросы ❔", buttons.FAQ)
}

func (h *handler) GetFAQAnswer(ctx context.Context, chatID int64, args []string) error {
	qIdxStr := args[0]
	qIdx, err := strconv.Atoi(qIdxStr)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "GetFAQAnswer",
			CausedBy:    "Atoi",
		})
	}
	ans := templates.GetQuestion(qIdx) + "\n\n" + templates.GetAnswer(qIdx)
	if err := h.sendMessage(chatID, ans); err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "GetFAQAnswer",
			CausedBy:    "sendMessage",
		})
	}
	return h.sendWithKeyboard(chatID, "Остались вопросы?", buttons.AskMore)
}
