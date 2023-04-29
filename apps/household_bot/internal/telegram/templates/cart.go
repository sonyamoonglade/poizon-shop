package templates

import (
	"fmt"
	"strings"

	"domain"
)

const (
	cartPreviewStart = "Вот твоя корзина!\nПозиций в корзине: %d\n\n"
	cartPreviewEnd   = "Итого в рублях: %d ₽\n\n" +
		"Ты выбрал товар(ы) под заказ, с учетом упаковки, страховки и доставки до Москвы, " +
		"дальнейшая отправка в другие города считается и оплачивается отдельно в ТК 🌉\n\n" +
		"Готов заказать? Жми кнопку!"

	cartPreviewDiscountEnd = "Итого в рублях (с учетом скидки): %d ₽\n\n" +
		"Ты выбрал товар(ы) под заказ, с учетом упаковки, страховки и доставки до Москвы, " +
		"дальнейшая отправка в другие города считается и оплачивается отдельно в ТК 🌉\n\n" +
		"Готов заказать? Жми кнопку!"

	positionPreviewInStock = "Количество: %d\nНазвание: %s\nЦена: %d ₽\nЦена по рынку: %d ₽\nНаличие: %s\nАртикул: %s\n\n"
	positionPreviewOrdered = "Количество: %d\nНазвание: %s\nЦена: %d ₽\nЦена по рынку: %d ₽\nАртикул: %s\n\n"

	positionPreviewWithDiscountInStock = "Количество: %d\nНазвание: %s\nЦена: %d ₽\nЦена (с учетом скидки): %d ₽\nЦена по рынку: %d ₽\nНаличие: %s\nАртикул: %s\n\n"
	positionPreviewWithDiscountOrdered = "Количество: %d\nНазвание: %s\nЦена: %d ₽\nЦена (с учетом скидки): %d ₽\nЦена по рынку: %d ₽\nАртикул: %s\n\n"

	editPosition = "Выбери номер позиции, чтобы удалить её 🙅‍♂️\n\nПо клику на " +
		"кнопку позиция изчезнет из твоей корзины!"

	tryToAddWithInvalidInStock = "Ты пытаешься добавить товар \"%s\", но в твоей корзине уже присутствует " +
		"товар \"%s\".\n\nОчисти корзину или перейди в каталог \"%s\"."

	checkingCart = "Проверяем твою корзину... Пожалуйста, подожди"

	productNotFound = "Продукт: %s с артикулом: [%s] Oтсутствует!\n\nУдали его из корзины"
)

func TryAddWithInvalidInStock(actual, want bool) string {
	actualStr, wantStr := domain.InStockToString(actual), domain.InStockToString(want)
	return fmt.Sprintf(tryToAddWithInvalidInStock, actualStr, wantStr, wantStr)
}

func EditCartPosition() string {
	return editPosition
}

func CheckingCart() string {
	return checkingCart
}

func ProductNotFound(name, isbn string) string {
	return fmt.Sprintf(productNotFound, name, isbn)
}

func RenderCart(cart domain.HouseholdCart, inStock bool) string {
	start := _cartPreviewStart(len(cart))
	var total uint32
	for _, groupedProduct := range cart.Group() {
		pos := groupedProduct.P
		if inStock {
			start += _cartPositionInStock(cartPositionInStockArgs{
				qty:         groupedProduct.Qty,
				price:       pos.Price,
				priceGlob:   pos.PriceGlob,
				availableIn: *pos.AvailableIn,
				name:        pos.Name,
				isbn:        pos.ISBN,
			})
		} else {
			start += _cartPositionOrdered(cartPositionOrderedArgs{
				qty:       groupedProduct.Qty,
				price:     pos.Price,
				priceGlob: pos.PriceGlob,
				name:      pos.Name,
				isbn:      pos.ISBN,
			})
		}

		total += pos.Price * uint32(groupedProduct.Qty)
	}
	return start + _cartPreviewEnd(total, false)
}

func RenderCartWithDiscount(cart domain.HouseholdCart, discount uint32, inStock bool) string {
	start := _cartPreviewStart(cart.Size())

	var discountedTotal uint32
	for _, groupedProduct := range cart.Group() {
		pos := groupedProduct.P
		if inStock {
			start += _cartPositionWithDiscountInStock(cartPositionWithDiscountInStockArgs{
				qty:             groupedProduct.Qty,
				price:           pos.Price,
				priceGlob:       pos.PriceGlob,
				availableIn:     *pos.AvailableIn,
				name:            pos.Name,
				isbn:            pos.ISBN,
				discountedPrice: pos.Price - discount,
			})
		} else {
			start += _cartPositionWithDiscountOrdered(cartPositionWithDiscountOrderedArgs{
				qty:             groupedProduct.Qty,
				price:           pos.Price,
				priceGlob:       pos.PriceGlob,
				name:            pos.Name,
				isbn:            pos.ISBN,
				discountedPrice: pos.Price - discount,
			})
		}

		discountedTotal += pos.Price*uint32(groupedProduct.Qty) - discount
	}
	return start + _cartPreviewEnd(discountedTotal, true)
}

type cartPositionOrderedArgs struct {
	qty       int
	price     uint32
	priceGlob uint32
	name      string
	isbn      string
}

func _cartPositionOrdered(args cartPositionOrderedArgs) string {
	return fmt.Sprintf(
		positionPreviewOrdered,
		args.qty,
		args.name,
		args.price,
		args.priceGlob,
		args.isbn,
	)
}

type cartPositionInStockArgs struct {
	qty         int
	price       uint32
	priceGlob   uint32
	name        string
	availableIn []string
	isbn        string
}

func _cartPositionInStock(args cartPositionInStockArgs) string {
	return fmt.Sprintf(
		positionPreviewInStock,
		args.qty,
		args.name,
		args.price,
		args.priceGlob,
		strings.Join(args.availableIn, ";"),
		args.isbn,
	)
}

type cartPositionWithDiscountOrderedArgs struct {
	qty             int
	price           uint32
	discountedPrice uint32
	priceGlob       uint32
	name            string
	isbn            string
}

func _cartPositionWithDiscountOrdered(args cartPositionWithDiscountOrderedArgs) string {
	return fmt.Sprintf(
		positionPreviewWithDiscountOrdered,
		args.qty,
		args.name,
		args.price,
		args.discountedPrice,
		args.priceGlob,
		args.isbn,
	)
}

type cartPositionWithDiscountInStockArgs struct {
	qty             int
	price           uint32
	discountedPrice uint32
	priceGlob       uint32
	availableIn     []string
	name            string
	isbn            string
}

func _cartPositionWithDiscountInStock(args cartPositionWithDiscountInStockArgs) string {
	return fmt.Sprintf(
		positionPreviewWithDiscountInStock,
		args.qty,
		args.name,
		args.price,
		args.discountedPrice,
		args.priceGlob,
		strings.Join(args.availableIn, ";"),
		args.isbn,
	)
}

func _cartPreviewEnd(totalRub uint32, discounted bool) string {
	if discounted {
		return fmt.Sprintf(cartPreviewDiscountEnd, totalRub)
	}
	return fmt.Sprintf(cartPreviewEnd, totalRub)
}

func _cartPreviewStart(numPositions int) string {
	return fmt.Sprintf(cartPreviewStart, numPositions)
}
