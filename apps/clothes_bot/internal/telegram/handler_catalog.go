package telegram

import (
	"context"
	"fmt"

	"domain"
	"dto"
	"functools"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/go-cmp/cmp"
)

func (h *handler) Catalog(ctx context.Context, chatID int64) error {
	var (
		telegramID = chatID
		first      bool
	)

	customer, err := h.customerService.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return err
	}

	if err := h.sendMessage(chatID, getCatalog(*customer.Username)); err != nil {
		return err
	}

	// Load appropriate item
	item := h.catalogProvider.LoadAt(customer.CatalogOffset)
	if cmp.Equal(item, domain.ClothingProduct{}) {
		return h.sendMessage(chatID, "Каталог отсутствует")
	}

	var (
		discount   = new(uint32)
		discounted bool
	)
	if customer.HasPromocode() {
		discounted = true
		*discount = customer.MustGetPromocode().GetClothingDiscount()
	}
	thumbnails := functools.Map(func(url string, i int) interface{} {
		thumbnail := tg.NewInputMediaPhoto(tg.FileURL(url))
		if !first {
			// add caption to first element
			thumbnail.Caption = productCard(item, discounted, discount)
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
		next, _ := h.catalogProvider.LoadNext(currentOffset)
		btnArgs.nextTitle = next.Title
	}

	if hasPrev {
		prev, _ := h.catalogProvider.LoadPrev(currentOffset)
		btnArgs.prevTitle = prev.Title
	}

	// Do not send buttons
	if !hasPrev && !hasNext {
		return nil
	}

	buttons := prepareCatalogButtons(btnArgs)
	return h.sendWithKeyboard(chatID, "Кнопки для пролистывания каталога", buttons)
}

// No need to call h.catalogProvider.HasNext. See h.Catalog impl
func (h *handler) HandleCatalogNext(ctx context.Context, chatID int64, controlButtonsMsgID int64, thumbnailMsgIDs []int) error {
	var telegramID = chatID
	customer, err := h.customerService.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return err
	}
	next, ok := h.catalogProvider.LoadNext(customer.CatalogOffset)
	if !ok {
		next, ok = h.catalogProvider.LoadFirst()
		customer.NullifyCatalogOffset()

		if !ok {
			return h.sendMessage(chatID, "Каталог отсутствует")
		}

		return h.sendMessage(chatID, "Товар был удален. Открой каталог заново")
	}
	// Increment the offset
	customer.CatalogOffset++
	return h.updateCatalog(ctx, chatID, thumbnailMsgIDs, controlButtonsMsgID, customer, next)
}

func (h *handler) HandleCatalogPrev(ctx context.Context, chatID int64, controlButtonsMsgID int64, thumbnailMsgIDs []int) error {
	var telegramID = chatID
	customer, err := h.customerService.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return err
	}
	prev, ok := h.catalogProvider.LoadPrev(customer.CatalogOffset)
	if !ok {
		prev, ok = h.catalogProvider.LoadFirst()
		customer.NullifyCatalogOffset()

		if !ok {
			return h.sendMessage(chatID, "Каталог отсутствует")
		}
		if err := h.sendMessage(chatID, "Видимо каталог обновили. Перемещаем тебя в начало!"); err != nil {
			return err
		}
		return h.updateCatalog(ctx, chatID, thumbnailMsgIDs, controlButtonsMsgID, customer, prev)
	}
	// Decrement the offset
	customer.CatalogOffset--
	return h.updateCatalog(ctx, chatID, thumbnailMsgIDs, controlButtonsMsgID, customer, prev)
}

func (h *handler) updateCatalog(ctx context.Context,
	chatID int64,
	thumbnailMsgIDs []int,
	controlButtonsMsgID int64,
	customer domain.ClothingCustomer,
	item domain.ClothingProduct) error {
	// Null
	if cmp.Equal(item, domain.ClothingProduct{}) {
		if err := h.deleteBulk(chatID, append(thumbnailMsgIDs, int(controlButtonsMsgID))...); err != nil {
			return fmt.Errorf("delete bulk: %w", err)
		}
		return h.sendWithKeyboard(chatID, "Данный товар больше не существует. Зайди в каталог заново", prepareMenuButtons(false))
	}
	var (
		discount   = new(uint32)
		discounted bool
	)
	if customer.HasPromocode() {
		discounted = true
		*discount = customer.MustGetPromocode().GetClothingDiscount()
	}
	// Get next item images
	var first bool
	thumbnails := functools.Map(func(url string, i int) interface{} {
		thumbnail := tg.NewInputMediaPhoto(tg.FileURL(url))
		if !first {
			// add caption to first element
			thumbnail.Caption = productCard(item, discounted, discount)
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
				MessageID: thumbMsgID,
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
	if err := h.customerService.Update(ctx, customer.CustomerID, updateDTO); err != nil {
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
		next, _ := h.catalogProvider.LoadNext(customer.CatalogOffset)
		btnArgs.nextTitle = next.Title
	}

	if hasPrev {
		prev, _ := h.catalogProvider.LoadPrev(customer.CatalogOffset)
		btnArgs.prevTitle = prev.Title
	}

	buttons := prepareCatalogButtons(btnArgs)
	editButtons := tg.NewEditMessageReplyMarkup(chatID, int(controlButtonsMsgID), buttons)

	return h.cleanSend(editButtons)
}
