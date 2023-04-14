package templates

import (
	"fmt"

	"domain"
)

const (
	cartPreviewStartTemplate = "Вот твоя корзина!\nПозиций в корзине: %d\n"
	cartPreviewEndTemplate   = "Итого в рублях: %d ₽\n\n" +
		"В стоимость каждой позиции включена страховка и доставка до Москвы\n\n---\n\nГотов заказать? Жми на кнопку!"
	positionPreviewTemplate = "%d. name: %s\nprice: %d\n"

	editPositionTemplate = "Выбери номер позиции, чтобы удалить её 🙅‍♂️\n\nПо клику на " +
		"кнопку позиция изчезнет из твоей корзины!"

	tryToAddWithInvalidInStockTemplate = "Вы пытаетесь добавить товар '%s', но в вашей корзине уже присутствует" +
		"товар '%s'.\n\nОчистите корзину или перейдите в каталог '%s'"
)

func TryAddWithInvalidInStock(actual, want bool) string {
	actualStr, wantStr := domain.InStockToString(actual), domain.InStockToString(want)
	return fmt.Sprintf(tryToAddWithInvalidInStockTemplate, actualStr, wantStr, wantStr)
}

type cartPositionArgs struct {
	n           int
	priceRub    uint32
	productName string
}

func EditCartPosition() string {
	return editPositionTemplate
}

func RenderCart(cart domain.HouseholdCart) string {
	start := _cartPreviewStartTemplate(len(cart))
	var total uint32
	for i, pos := range cart {
		start += _cartPositionTemplate(cartPositionArgs{
			n:           i + 1,
			priceRub:    pos.Price,
			productName: pos.Name,
		})
		total += pos.Price
	}
	return start + _cartPreviewEndTemplate(total)
}

func _cartPositionTemplate(args cartPositionArgs) string {
	return fmt.Sprintf(positionPreviewTemplate, args.n, args.productName, args.priceRub)
}

func _cartPreviewEndTemplate(totalRub uint32) string {
	return fmt.Sprintf(cartPreviewEndTemplate, totalRub)
}

func _cartPreviewStartTemplate(numPositions int) string {
	return fmt.Sprintf(cartPreviewStartTemplate, numPositions)
}
