package telegram

import (
	"context"
	"reflect"

	"domain"
	"dto"
	"functools"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *handler) Catalog(ctx context.Context, chatID int64) error {
	var (
		telegramID = chatID
		first      bool
	)

	customer, err := h.customerRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return err
	}

	if err := h.sendMessage(chatID, getCatalog(*customer.Username)); err != nil {
		return err
	}

	// Load appropriate item
	item := h.catalogProvider.LoadAt(customer.CatalogOffset)

	thumbnails := functools.Map(func(url string, i int) interface{} {
		thumbnail := tg.NewInputMediaPhoto(tg.FileURL(url))
		if !first {
			// add caption to first element
			thumbnail.Caption = item.GetCaption()
			first = true
		}
		thumbnail.ParseMode = parseModeHTML
		return thumbnail
	}, item.ImageURLs)

	group := tg.NewMediaGroup(chatID, thumbnails)

	// Sends thumnails with caption
	sentMsgs, err := h.b.SendMediaGroup(group)
	if err != nil {
		return err
	}

	msgIDs := functools.Map(func(m tg.Message, i int) int {
		return m.MessageID
	}, sentMsgs)

	// Prepare buttons for controlling prev, next
	var (
		currentOffset = customer.CatalogOffset
		hasNext       = h.catalogProvider.HasNext(currentOffset)
		hasPrev       = h.catalogProvider.HasPrev(currentOffset)
	)

	btnArgs := catalogButtonsArgs{
		hasNext: hasNext,
		hasPrev: hasPrev,
		msgIDs:  msgIDs,
	}

	if hasNext {
		next := h.catalogProvider.LoadNext(currentOffset)
		btnArgs.nextTitle = next.Title
	}

	if hasPrev {
		prev := h.catalogProvider.LoadPrev(currentOffset)
		btnArgs.prevTitle = prev.Title
	}

	if !hasPrev && !hasNext {
		return nil
	}

	buttons := prepareCatalogButtons(btnArgs)
	return h.sendWithKeyboard(chatID, "Кнопки для пролистывания каталога", buttons)
}

// No need to call h.catalogProvider.HasNext. See h.Catalog impl
func (h *handler) HandleCatalogNext(ctx context.Context, chatID int64, controlButtonsMsgID int64, thumbnailMsgIDs []int) error {
	var telegramID = chatID
	customer, err := h.customerRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return err
	}

	next := h.catalogProvider.LoadNext(customer.CatalogOffset)
	// Increment the offset
	customer.CatalogOffset++

	return h.updateCatalog(ctx, chatID, thumbnailMsgIDs, controlButtonsMsgID, customer, next)
}

func (h *handler) HandleCatalogPrev(ctx context.Context, chatID int64, controlButtonsMsgID int64, thumbnailMsgIDs []int) error {
	var telegramID = chatID
	customer, err := h.customerRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return err
	}

	prev := h.catalogProvider.LoadPrev(customer.CatalogOffset)
	// Decrement the offset
	customer.CatalogOffset--

	return h.updateCatalog(ctx, chatID, thumbnailMsgIDs, controlButtonsMsgID, customer, prev)
}

func (h *handler) updateCatalog(ctx context.Context,
	chatID int64,
	thumbnailMsgIDs []int,
	controlButtonsMsgID int64,
	customer domain.Customer,
	item domain.ClothingProduct) error {
	// Null
	if reflect.DeepEqual(domain.ClothingProduct{}, item) {
		return nil
	}
	// Get next item images
	var first bool
	thumbnails := functools.Map(func(url string, i int) interface{} {
		thumbnail := tg.NewInputMediaPhoto(tg.FileURL(url))
		if !first {
			// add caption to first element
			thumbnail.Caption = item.GetCaption()
			first = true
		}
		thumbnail.ParseMode = parseModeHTML
		return thumbnail
	}, item.ImageURLs)

	var sentMsgIDs []int
	// Draw it by updating
	for i, thumbMsgID := range thumbnailMsgIDs {
		editOneMedia := &tg.EditMessageMediaConfig{
			BaseEdit: tg.BaseEdit{
				ChatID:    chatID,
				MessageID: int(thumbMsgID),
			},
			Media: thumbnails[i],
		}
		sentMessage, err := h.b.Send(editOneMedia)
		if err != nil {
			return err
		}
		sentMsgIDs = append(sentMsgIDs, sentMessage.MessageID)
	}

	updateDTO := dto.UpdateClothingCustomerDTO{
		CatalogOffset: &customer.CatalogOffset,
	}
	if err := h.customerRepo.Update(ctx, customer.CustomerID, updateDTO); err != nil {
		return err
	}

	// Load accordingly to next offset
	var (
		hasNext = h.catalogProvider.HasNext(customer.CatalogOffset)
		hasPrev = h.catalogProvider.HasPrev(customer.CatalogOffset)
	)

	// update buttons
	btnArgs := catalogButtonsArgs{
		hasNext: hasNext,
		hasPrev: hasPrev,
		msgIDs:  sentMsgIDs,
	}
	if hasNext {
		next := h.catalogProvider.LoadNext(customer.CatalogOffset)
		btnArgs.nextTitle = next.Title
	}

	if hasPrev {
		prev := h.catalogProvider.LoadPrev(customer.CatalogOffset)
		btnArgs.prevTitle = prev.Title
	}

	buttons := prepareCatalogButtons(btnArgs)
	editButtons := tg.NewEditMessageReplyMarkup(chatID, int(controlButtonsMsgID), buttons)

	return h.cleanSend(editButtons)
}
