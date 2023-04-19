package templates

import (
	"fmt"

	"domain"
)

const (
	cartPreviewStartTemplate = "Вот твоя корзина!\nПозиций в корзине: %d\n\n"
	cartPreviewEndTemplate   = "Итого в рублях: %d ₽\n\n" +
		"Ты выбрал товар(ы) под заказ, с учетом упаковки, страховки и доставки до Москвы, " +
		"дальнейшая отправка в другие города считается и оплачивается отдельно в ТК 🌉\n\n" +
		"Готов заказать? Жми кнопку!"
	positionPreviewTemplate = "%d. Название: %s\nЦена: %d ₽\nЦена по рынку: %d ₽\nНаличие: [todo?]\nАртикул: %s\n\n"

	editPositionTemplate = "Выбери номер позиции, чтобы удалить её 🙅‍♂️\n\nПо клику на " +
		"кнопку позиция изчезнет из твоей корзины!"

	tryToAddWithInvalidInStockTemplate = "Ты пытаешься добавить товар \"%s\", но в твоей корзине уже присутствует " +
		"товар \"%s\".\n\nОчисти корзину или перейди в каталог \"%s\"."
)

func TryAddWithInvalidInStock(actual, want bool) string {
	actualStr, wantStr := domain.InStockToString(actual), domain.InStockToString(want)
	return fmt.Sprintf(tryToAddWithInvalidInStockTemplate, actualStr, wantStr, wantStr)
}

type cartPositionArgs struct {
	n         int
	price     uint32
	priceGlob uint32
	name      string
	isbn      string
}

func EditCartPosition() string {
	return editPositionTemplate
}

func RenderCart(cart domain.HouseholdCart) string {
	start := _cartPreviewStartTemplate(len(cart))
	var total uint32
	for i, pos := range cart {
		start += _cartPositionTemplate(cartPositionArgs{
			n:         i + 1,
			price:     pos.Price,
			priceGlob: pos.PriceGlob,
			name:      pos.Name,
			isbn:      pos.ISBN,
		})
		total += pos.Price
	}
	return start + _cartPreviewEndTemplate(total)
}

func _cartPositionTemplate(args cartPositionArgs) string {
	return fmt.Sprintf(
		positionPreviewTemplate,
		args.n,
		args.name,
		args.price,
		args.priceGlob,
		args.isbn,
	)
}

func _cartPreviewEndTemplate(totalRub uint32) string {
	return fmt.Sprintf(cartPreviewEndTemplate, totalRub)
}

func _cartPreviewStartTemplate(numPositions int) string {
	return fmt.Sprintf(cartPreviewStartTemplate, numPositions)
}
