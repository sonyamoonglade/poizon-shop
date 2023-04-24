package templates

import (
	"fmt"

	"domain"
)

const (
	askForDeliveryAddress = "Отправь адрес ближайшего постамата PickPoint или отделения СДЭК ⛳️ в формате:\n\n" +
		"Страна, область, город, улица, номер дома/строения 🏡\n\n" +
		"Я доставлю твой заказ туда 🚚"

	askForPhoneNumber = "Отправь мне свой контактный номер телефона в формате:\n 👉 79128000000"

	invalidFIOInput = "Неправильный формат полного имени.\nОтправь полное имя в " +
		"формате - Иванов Иван Иванович"

	askForFIO = "Укажи ФИО получателя \U0001FAAA"

	deliveryOnlyToMoscow = "Стоимость указана с учетом доставки товара из Китая до Москвы, доставка в другие " +
		"города и районы России просчитывается и оплачивается отдельно в ТК СДЕК 🚚"

	requisites = "Счет для оплаты заказа: [%s]\n\nТимофеев Вадим Денисович 🙋🏼‍♂️️@xKK_Russia\n\nНомер карты " +
		"Сбер: %s\nНомер карты Тинькофф: %s\nВ комментарии укажи номер заказа [%s]\n\nПосле оплаты нажми кнопку «Оплачено»\n"

	cartWarn = "Вся информация, которую ты указываешь, для сборки в корзине 🧺 должна быть актуальной, " +
		"если она составляет более 48ч ⌚️и является неактуальной  – заказ не будет принят и деньги возвратятся в полном " +
		"объеме на карту плательщика 💴\n\n"

	orderStart = "%s - Твоя заявка готова!\n\nНомер заказа: [%s]\n\nДанные " +
		"получателя:\nФИО: %s\nНомер телефона: %s\nАдрес доставки: %s\n\nТоваров в корзине: %d\n\n"

	orderEnd = "Итоговая стоимость составляет %d ₽\n"

	orderEndDiscount = "Итоговая стоимость с учетом скидки составляет %d ₽\n"

	sendingRequsities = "\nВысылаю реквизиты для оплаты 🧾"

	successfulPayment = "%s, твой заказ %s сейчас на подтверждении у админа. Он напишет тебе в личные сообщения " +
		"и подтвердит статус покупки.\n\n‼️Никому кроме бота деньги отправлять не нужно‼️Даже админу‼️\n\nТолько " +
		"админ проверяет поступление денег и обозначает статус покупки ✅"

	myOrdersBody = "Заказ: %s\nАдрес доставки: %s\n\nОплачен: %s\nПодтвержден админом: %s\n" +
		"Статус заказа: %s\n\nТоваров в корзине: %d\nСумма в рублях: %d ₽\n\n" +
		"Комментарий админа: %s\n\nТовар(ы):\n"

	myOrdersStart     = "Вот твои заказы, %s!\n\n"
	myOrdersSeparator = "\n-----\n\n"

	yes string = "✅"
	no         = "❌"
)

func AskForFIO() string {
	return askForFIO
}

func InvalidFIO() string {
	return invalidFIOInput
}

func AskForPhoneNumber() string {
	return askForPhoneNumber
}

func AskForDeliveryAddress() string {
	return askForDeliveryAddress
}

func Requisites(shortOrderID string, r domain.Requisites) string {
	return fmt.Sprintf(requisites, shortOrderID, r.SberID, r.TinkoffID, shortOrderID)
}

func RenderOrderAfterPayment(order domain.HouseholdOrder) string {
	start := cartWarn + _orderStart(order)
	for i, pos := range order.Cart {
		if order.InStock {
			start += _cartPositionInStock(cartPositionInStockArgs{
				n:           i + 1,
				price:       pos.Price,
				availableIn: *pos.AvailableIn,
				priceGlob:   pos.PriceGlob,
				name:        pos.Name,
				isbn:        pos.ISBN,
			})
		} else {
			start += _cartPositionOrdered(cartPositionOrderedArgs{
				n:         i + 1,
				price:     pos.Price,
				priceGlob: pos.PriceGlob,
				name:      pos.Name,
				isbn:      pos.ISBN,
			})
		}

	}
	return start + _orderEnd(order.AmountRUB) + sendingRequsities
}

func RenderOrderAfterPaymentWithDiscount(order domain.HouseholdOrder, discount uint32) string {
	start := cartWarn + _orderStart(order)
	for i, pos := range order.Cart {

		if order.InStock {
			start += _cartPositionWithDiscountInStock(cartPositionWithDiscountInStockArgs{
				n:               i + 1,
				price:           pos.Price,
				discountedPrice: pos.Price - discount,
				availableIn:     *pos.AvailableIn,
				priceGlob:       pos.PriceGlob,
				name:            pos.Name,
				isbn:            pos.ISBN,
			})
		} else {
			start += _cartPositionWithDiscountOrdered(cartPositionWithDiscountOrderedArgs{
				n:               i + 1,
				price:           pos.Price,
				discountedPrice: pos.Price - discount,
				priceGlob:       pos.PriceGlob,
				name:            pos.Name,
				isbn:            pos.ISBN,
			})
		}

	}
	return start + _orderEndWithDiscount(order.DiscountedAmount) + sendingRequsities
}

func RenderMyOrders(name string, orders []domain.HouseholdOrder) string {
	start := fmt.Sprintf(myOrdersStart, name)
	for _, o := range orders {
		start += _myOrderBody(o)
		for i, pos := range o.Cart {
			if o.InStock {
				start += _cartPositionInStock(cartPositionInStockArgs{
					n:           i + 1,
					price:       pos.Price,
					priceGlob:   pos.PriceGlob,
					availableIn: *pos.AvailableIn,
					name:        pos.Name,
					isbn:        pos.ISBN,
				})
			} else {
				start += _cartPositionOrdered(cartPositionOrderedArgs{
					n:         i + 1,
					price:     pos.Price,
					priceGlob: pos.PriceGlob,
					name:      pos.Name,
					isbn:      pos.ISBN,
				})
			}
		}
		start += _orderEnd(o.AmountRUB) + myOrdersSeparator
	}
	return start
}

func SuccessfulPayment(fullName, orderShortID string) string {
	return fmt.Sprintf(successfulPayment, fullName, orderShortID)
}

func _orderStart(order domain.HouseholdOrder) string {
	return fmt.Sprintf(
		orderStart,
		*order.Customer.FullName,
		order.ShortID,
		*order.Customer.FullName,
		*order.Customer.PhoneNumber,
		order.DeliveryAddress,
		len(order.Cart),
	)
}

func _orderEnd(totalRub uint32) string {
	return fmt.Sprintf(orderEnd, totalRub)
}
func _orderEndWithDiscount(totalRub uint32) string {
	return fmt.Sprintf(orderEndDiscount, totalRub)
}

func _myOrderBody(o domain.HouseholdOrder) string {
	return fmt.Sprintf(
		myOrdersBody,
		o.ShortID,
		o.DeliveryAddress,
		formatBool(o.IsPaid),
		formatBool(o.IsApproved),
		domain.StatusTexts[o.Status],
		len(o.Cart),
		o.AmountRUB,
		o.GetComment(),
	)
}

func formatBool(b bool) string {
	if b {
		return yes
	}
	return no
}
