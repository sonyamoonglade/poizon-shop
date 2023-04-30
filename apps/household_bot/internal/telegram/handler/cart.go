package handler

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"household_bot/internal/telegram/callback"
	"household_bot/internal/telegram/router"

	fn "github.com/sonyamoonglade/go_func"

	"domain"
	"dto"
	"household_bot/internal/telegram/buttons"
	"household_bot/internal/telegram/templates"
	"household_bot/internal/telegram/tg_errors"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *handler) AddToCart(ctx context.Context, chatID int64, prevMsgID int, args []string, source router.AddToCartSource) error {
	var telegramID = chatID

	customer, err := h.customerService.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "AddToCart",
			CausedBy:    "GetByTelegramID",
		})
	}

	var currentProduct domain.HouseholdProduct
	if source == router.SourceCatalog && args != nil && len(args) > 2 {
		cp, err := h.addToCartFromCatalog(ctx, args)
		if err != nil {
			return tg_errors.New(tg_errors.Config{
				OriginalErr: err,
				Handler:     "AddToCart",
				CausedBy:    "addToCartFromCatalog",
			})
		}
		currentProduct = cp
	} else if source == router.SourceISBNSearch {
		cp, err := h.addToCartFromISBNSearch(args)
		if err != nil {
			return tg_errors.New(tg_errors.Config{
				OriginalErr: err,
				Handler:     "AddToCart",
				CausedBy:    "addToCartFromISBNSearch",
			})
		}
		currentProduct = cp
	}

	// Category that customer is adding product with
	currentCategory, _ := h.catalogProvider.GetCategoryByID(currentProduct.CategoryID)

	firstProduct, exists := customer.Cart.First()

	if customer.Cart.IsEmpty() {
		customer.Cart.Add(currentProduct)
	} else if exists {
		// Have to check if currentCategory is the same as 0th element in cart
		firstProductCategory, err := h.categoryService.GetByID(ctx, firstProduct.CategoryID)
		if err != nil {
			if errors.Is(err, domain.ErrCategoryNotFound) {
				return h.handleIfCategoryNotFound(ctx, chatID, customer, firstProduct)
			}
			return tg_errors.New(tg_errors.Config{
				OriginalErr: err,
				Handler:     "AddToCart",
				CausedBy:    "GetByID",
			})
		}
		// If inStock field is the same then it's fine to add
		if firstProductCategory.InStock == currentCategory.InStock {
			customer.Cart.Add(currentProduct)
		} else {
			return h.sendWithKeyboard(
				chatID,
				templates.TryAddWithInvalidInStock(
					currentCategory.InStock,
					!currentCategory.InStock,
				),
				buttons.RouteToCatalogOrCart,
			)
		}
	}

	err = h.customerService.Update(ctx, customer.CustomerID, dto.UpdateHouseholdCustomerDTO{
		Cart: &customer.Cart,
	})
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "AddToCart",
			CausedBy:    "Update",
		})
	}

	var (
		cTitle,
		sTitle,
		inStockStr,
		pName = args[0], args[1], args[2], args[3]
	)
	inStock, err := strconv.ParseBool(inStockStr)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "AddToCart",
			CausedBy:    "ParseBool",
		})
	}
	currentQuantity := fn.Reduce(func(acc int, el domain.HouseholdProduct, _ int) int {
		if pName == el.Name {
			acc += 1
		}
		return acc
	}, customer.Cart.Slice(), 0)
	keyboard := buttons.NewProductCardButtons(buttons.ProductCardButtonsArgs{
		Cb:       callback.AddToCart,
		CTitle:   cTitle,
		STitle:   sTitle,
		PName:    pName,
		InStock:  inStock,
		Quantity: currentQuantity,
		Back: buttons.NewBackButton(callback.FromProductCardToProducts,
			&cTitle,
			&sTitle,
			&inStock,
		),
	})

	editMsg := tg.NewEditMessageReplyMarkup(chatID, prevMsgID, keyboard)
	return h.cleanSend(editMsg)
	//return h.sendWithKeyboard(chatID, templates.PositionAdded(currentProduct.Name), buttons.GotoCart)
}

func (h *handler) GetCart(ctx context.Context, chatID int64) error {
	var telegramID = chatID
	customer, err := h.customerService.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "AddToCart",
			CausedBy:    "GetByTelegramID",
		})
	}
	if customer.Cart.IsEmpty() {
		return h.sendWithKeyboard(chatID, "Твоя корзина пуста!\n\nДобавим что-нибудь?", buttons.AddPosition)
	}

	firstProduct, _ := customer.Cart.First()
	category, _ := h.catalogProvider.GetCategoryByID(firstProduct.CategoryID)

	if customer.HasPromocode() {
		promo := customer.MustGetPromocode()
		return h.sendWithKeyboard(
			chatID,
			templates.RenderCartWithDiscount(
				customer.Cart,
				promo.GetHouseholdDiscount(),
				category.InStock,
			),
			buttons.CartPreview,
		)
	}

	return h.sendWithKeyboard(chatID, templates.RenderCart(customer.Cart, category.InStock), buttons.CartPreview)
}

func (h *handler) EditCart(ctx context.Context, chatID int64, cartMsgID int) error {
	var telegramID = chatID

	customer, err := h.customerService.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return fmt.Errorf("customerService.GetByTelegramID: %w", err)
	}

	if customer.Cart.IsEmpty() {
		return h.emptyCart(chatID)
	}

	if err := h.customerService.UpdateState(ctx, telegramID, domain.StateWaitingForCartPositionToEdit); err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "EditCart",
			CausedBy:    "UpdateState",
		})
	}

	return h.sendWithKeyboard(chatID,
		templates.EditCartPosition(),
		buttons.NewEditCartButtonsGroup(
			customer.Cart.Group(),
			cartMsgID,
		),
	)
}

func (h *handler) DeletePositionFromCart(ctx context.Context, chatID int64, buttonsMsgID int, args []string) error {
	var (
		cartMsgIDStr,
		productIDStr = args[0], args[1]
	)

	cartMsgID, err := strconv.Atoi(cartMsgIDStr)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "RemoveCartPosition",
			CausedBy:    "Atoi",
		})
	}

	customer, err := h.checkRequiredState(ctx, chatID, domain.StateWaitingForCartPositionToEdit)
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "RemoveCartPosition",
			CausedBy:    "checkRequiredState",
		})
	}
	customer.Cart.Remove(domain.RemoveByProductID(productIDStr))
	updateDTO := dto.UpdateHouseholdCustomerDTO{
		Cart: &customer.Cart,
	}
	if err := h.customerService.Update(ctx, customer.CustomerID, updateDTO); err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "RemoveCartPosition",
			CausedBy:    "Update",
		})
	}

	// If customer has emptied cart just now
	if customer.Cart.IsEmpty() {
		// Delete edit buttons
		if err := h.bot.CleanRequest(tg.NewDeleteMessage(chatID, buttonsMsgID)); err != nil {
			return tg_errors.New(tg_errors.Config{
				OriginalErr: err,
				Handler:     "RemoveCartPosition",
				CausedBy:    "CleanRequest",
			})
		}
		// Update cart message
		msg := tg.NewEditMessageText(chatID, cartMsgID, "Ваша корзина пуста!")
		keyboard := buttons.AddPosition
		msg.ReplyMarkup = &keyboard
		if err := h.cleanSend(msg); err != nil {
			return tg_errors.New(tg_errors.Config{
				OriginalErr: err,
				Handler:     "RemoveCartPosition",
				CausedBy:    "cleanSend",
			})
		}
		return nil
	}

	// Edit original preview cart message and edit buttons
	firstProduct, _ := customer.Cart.First()
	category, _ := h.catalogProvider.GetCategoryByID(firstProduct.CategoryID)

	cartMsg := tg.NewEditMessageText(
		chatID,
		cartMsgID,
		templates.RenderCart(customer.Cart, category.InStock),
	)
	cartMsg.ReplyMarkup = &buttons.CartPreview

	updateButtons := tg.NewEditMessageReplyMarkup(chatID,
		buttonsMsgID,
		buttons.NewEditCartButtonsGroup(
			customer.Cart.Group(),
			cartMsgID,
		),
	)

	if err := h.sendBulk(updateButtons, cartMsg); err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "RemoveCartPosition",
			CausedBy:    "sendBulk",
		})
	}
	return nil
}

func (h *handler) emptyCart(chatID int64) error {
	return h.sendWithKeyboard(chatID, "Твоя корзина пуста!", buttons.AddPosition)
}

func (h *handler) addToCartFromCatalog(ctx context.Context, args []string) (domain.HouseholdProduct, error) {
	var (
		cTitle,
		sTitle,
		inStockStr,
		pName = args[0], args[1], args[2], args[3]
	)
	inStock, err := strconv.ParseBool(inStockStr)
	if err != nil {
		return domain.HouseholdProduct{}, tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "addToCartFromCatalog",
			CausedBy:    "ParseBool",
		})
	}

	return h.catalogProvider.GetProductByCategoryAndSubcategory(cTitle, sTitle, pName, inStock), nil
}

func (h *handler) addToCartFromISBNSearch(args []string) (domain.HouseholdProduct, error) {
	isbn := args[0]
	product, ok := h.catalogProvider.GetProductByISBN(isbn)
	if !ok {
		return domain.HouseholdProduct{}, domain.ErrProductNotFound
	}
	return product, nil
}

func (h *handler) handleIfCategoryNotFound(ctx context.Context, chatID int64, customer domain.HouseholdCustomer, product domain.HouseholdProduct) error {
	text := fmt.Sprintf("Категория с товаром (%s, %s) не найдена. ",
		product.Name,
		product.ISBN,
	)
	if err := h.sendMessage(chatID, text); err != nil {
		return err
	}

	if err := h.sendMessage(chatID, "Удаляем товар"); err != nil {
		return err
	}

	customer.Cart.RemoveAt(fn.
		Of(customer.Cart).
		IndexOf(func(cartProduct domain.HouseholdProduct) bool {
			return cartProduct.ProductID == product.ProductID
		}))

	err := h.customerService.Update(ctx, customer.CustomerID, dto.UpdateHouseholdCustomerDTO{
		Cart: &customer.Cart,
	})
	if err != nil {
		return tg_errors.New(tg_errors.Config{
			OriginalErr: err,
			Handler:     "handleIfCategoryNotFound",
			CausedBy:    "Update",
		})
	}
	return nil
}
