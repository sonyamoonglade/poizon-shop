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
)

func _cartPreviewStartTemplate(numPositions int) string {

	return fmt.Sprintf(cartPreviewStartTemplate, numPositions)
}

type cartPositionArgs struct {
	n           int
	priceRub    uint32
	productName string
}

func _cartPositionTemplate(args cartPositionArgs) string {
	return fmt.Sprintf(positionPreviewTemplate, args.n, args.productName, args.priceRub)
}

func _cartPreviewEndTemplate(totalRub uint32) string {
	return fmt.Sprintf(cartPreviewEndTemplate, totalRub)
}

type RenderCartArgs struct {
	Cart domain.HouseholdCart
}

func RenderCart(args RenderCartArgs) string {
	start := _cartPreviewStartTemplate(len(args.Cart))
	var total uint32
	for i, pos := range args.Cart {
		start += _cartPositionTemplate(cartPositionArgs{
			n:           i + 1,
			priceRub:    pos.Price,
			productName: pos.Name,
		})
		total += pos.Price
	}
	return start + _cartPreviewEndTemplate(total)
}
