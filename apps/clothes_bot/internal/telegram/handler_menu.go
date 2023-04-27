package telegram

import (
	"context"
	"errors"
	"fmt"

	"domain"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *handler) Start(ctx context.Context, m *tg.Message) error {
	var (
		chatID       = m.Chat.ID
		telegramID   = chatID
		chatUsername = m.From.String()
		username     = domain.MakeUsername(chatUsername)
	)
	// register customer
	_, err := h.customerService.GetByTelegramID(ctx, telegramID)
	if err != nil {
		if !errors.Is(err, domain.ErrCustomerNotFound) {
			return err
		}
		// save to db
		if err := h.customerService.Save(ctx, domain.NewClothingCustomer(telegramID, username)); err != nil {
			return err
		}
	}

	return h.sendWithKeyboard(chatID, getStartTemplate(username), initialMenuKeyboard)

}

func (h *handler) Menu(ctx context.Context, chatID int64) error {
	customer, err := h.customerService.GetByTelegramID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get by telegram id: %w", err)
	}
	if err := h.customerService.UpdateState(ctx, chatID, domain.StateDefault); err != nil {
		return fmt.Errorf("update state: %w", err)
	}
	rate, err := h.rateProvider.GetYuanRate(ctx)
	if err != nil {
		return fmt.Errorf("get yuan rate: %w", err)
	}
	if err := h.sendMessage(chatID, fmt.Sprintf("Курс юаня на сегодня: %.2f ₽", rate)); err != nil {
		return err
	}
	showPromoButton := !customer.HasPromocode()
	return h.sendWithKeyboard(chatID, getTemplate().Menu, prepareMenuButtons(showPromoButton))
}

func (h *handler) MyOrders(ctx context.Context, chatID int64) error {
	var telegramID = chatID

	customer, err := h.customerService.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return err
	}

	orders, err := h.orderService.GetAllForCustomer(ctx, customer.CustomerID)
	if err != nil {
		if errors.Is(err, domain.ErrNoOrders) {
			return h.sendMessage(chatID, "У тебя еще нет заказов 🦕")
		}

		return err
	}
	if len(orders) == 0 {
		return h.sendMessage(chatID, "У тебя еще нет заказов 🦕")
	}
	var name string
	if customer.FullName != nil {
		name = *customer.FullName
	} else {
		name = *customer.Username
	}
	out := getMyOrdersStart(name)

	var (
		discount   = new(uint32)
		discounted bool
	)
	if customer.HasPromocode() {
		discounted = true
		*discount = customer.MustGetPromocode().GetClothingDiscount()
	}
	for _, o := range orders {
		out += getSingleOrderPreview(o, customer.HasPromocode())
		for nCartItem, cartItem := range o.Cart {
			if discounted {
				out += getDiscountedPositionTemplate(cartPositionPreviewDiscountedArgs{
					n:           nCartItem + 1,
					link:        cartItem.ShopLink,
					size:        cartItem.Size,
					discountRub: *discount,
					category:    string(cartItem.Category),
					priceRub:    cartItem.PriceRUB,
					priceYuan:   cartItem.PriceYUAN,
				})
				continue
			}
			out += getPositionTemplate(cartPositionPreviewArgs{
				n:         nCartItem + 1,
				link:      cartItem.ShopLink,
				size:      cartItem.Size,
				category:  string(cartItem.Category),
				priceRub:  cartItem.PriceRUB,
				priceYuan: cartItem.PriceYUAN,
			})
		}

		out += getTemplate().MyOrdersEnd
	}

	return h.sendMessage(chatID, out)
}
