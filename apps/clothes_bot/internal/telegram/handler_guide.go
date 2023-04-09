package telegram

import (
	"context"

	"domain"
	"functools"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *handler) StartMakeOrderGuide(ctx context.Context, m *tg.Message) error {
	var (
		chatID     = m.Chat.ID
		telegramID = chatID
	)

	thumbnails := makeThumbnails(getTemplate().GuideStep1, guideStep1Thumbnail1, guideStep1Thumbnail2)
	group := tg.NewMediaGroup(chatID, thumbnails)

	sentMsgs, err := h.b.SendMediaGroup(group)
	if err != nil {
		return err
	}

	msgIDs := functools.Map(func(m tg.Message, i int) int {
		return m.MessageID
	}, sentMsgs)

	buttons := prepareOrderGuideButtons(orderGuideStep0Callback, msgIDs...)
	if err := h.sendWithKeyboard(chatID, "Кнопки для пролистывания инструкции", buttons); err != nil {
		return err
	}

	customer, err := h.customerRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return err
	}

	// If cart is not empty then skip order type ask
	if len(customer.Cart) > 0 {
		return h.askForCategory(ctx, chatID)
	}

	if err := h.customerRepo.UpdateState(ctx, telegramID, domain.StateWaitingForOrderType); err != nil {
		return err
	}

	return h.askForOrderType(ctx, chatID)
}

func (h *handler) MakeOrderGuideStep1(ctx context.Context, chatID int64, controlButtonsMessageID int, guideMsgIDs []int) error {
	thumbnails := makeThumbnails(getTemplate().GuideStep1, guideStep1Thumbnail1, guideStep1Thumbnail2)
	return h.updateGuideStep(chatID, guideMsgIDs, controlButtonsMessageID, orderGuideStep0Callback, thumbnails)
}

func (h *handler) MakeOrderGuideStep2(ctx context.Context, chatID int64, controlButtonsMessageID int, guideMsgIDs []int) error {
	thumbnails := makeThumbnails(getTemplate().GuideStep2, guideStep2Thumbnail1, guideStep2Thumbnail2)
	return h.updateGuideStep(chatID, guideMsgIDs, controlButtonsMessageID, orderGuideStep1Callback, thumbnails)
}

func (h *handler) MakeOrderGuideStep3(ctx context.Context, chatID int64, controlButtonsMessageID int, guideMsgIDs []int) error {
	thumbnails := makeThumbnails(getTemplate().GuideStep3, guideStep3Thumbnail1, guideStep3Thumbnail2)
	return h.updateGuideStep(chatID, guideMsgIDs, controlButtonsMessageID, orderGuideStep2Callback, thumbnails)
}

func (h *handler) MakeOrderGuideStep4(ctx context.Context, chatID int64, controlButtonsMessageID int, guideMsgIDs []int) error {
	thumbnails := makeThumbnails(getTemplate().GuideStep4, guideStep4Thumbnail1, guideStep4Thumbnail2)
	return h.updateGuideStep(chatID, guideMsgIDs, controlButtonsMessageID, orderGuideStep3Callback, thumbnails)
}

func (h *handler) MakeOrderGuideStep5(ctx context.Context, chatID int64, controlButtonsMessageID int, guideMsgIDs []int) error {
	thumbnails := makeThumbnails(getTemplate().GuideStep5, guideStep5Thumbnail1, guideStep5Thumbnail2)
	return h.updateGuideStep(chatID, guideMsgIDs, controlButtonsMessageID, orderGuideStep4Callback, thumbnails)
}

func (h *handler) MakeOrderGuideStep6(ctx context.Context, chatID int64, controlButtonsMessageID int, guideMsgIDs []int) error {
	thumbnails := makeThumbnails(getTemplate().GuideStep6, guideStep6Thumbnail1, guideStep6Thumbnail2)
	return h.updateGuideStep(chatID, guideMsgIDs, controlButtonsMessageID, orderGuideStep5Callback, thumbnails)
}

func (h *handler) updateGuideStep(chatID int64, guideMsgIDs []int, controlButtonsMessageID int, nextCallback int, thumbnails []interface{}) error {
	for i, t := range thumbnails {
		editOneMedia := &tg.EditMessageMediaConfig{
			BaseEdit: tg.BaseEdit{
				ChatID:    chatID,
				MessageID: guideMsgIDs[i],
			},
			Media: t,
		}
		if err := h.cleanSend(editOneMedia); err != nil {
			return err
		}
	}

	// update control buttons
	buttons := tg.NewEditMessageReplyMarkup(chatID, controlButtonsMessageID, prepareOrderGuideButtons(nextCallback, guideMsgIDs...))
	return h.cleanSend(buttons)

}
