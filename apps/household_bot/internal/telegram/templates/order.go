package templates

import (
	"fmt"

	"domain"
)

const (
	askForOrderTypeTemplate = "viberi tip zakaza"

	askForDeliveryAddressTemplate = "Отправь адрес ближайшего постамата PickPoint или отделения СДЭК ⛳️ в формате:\n\n" +
		"Страна, область, город, улица, номер дома/строения 🏡\n\n" +
		"Я доставлю твой заказ туда 🚚"

	askForPhoneNumberTemplate = "Отправь мне свой контактный номер телефона в формате:\n 👉 79128000000"

	invalidFIOInputTemplate = "Неправильный формат полного имени.\\n Отправь полное имя в " +
		"формате - Иванов Иван Иванович"

	askForFIOTemplate = "Укажи ФИО получателя \U0001FAAA"

	deliveryOnlyToMoscowTemplate = "Стоимость указана с учетом доставки товара из Китая до Москвы, доставка в другие " +
		"города и районы России просчитывается и оплачивается отдельно в ТК СДЕК 🚚"

	requisitesTemplate = "Счет для оплаты заказа: [%s]\n\nТимофеев Вадим Денисович 💁‍ ♂️ @xKK_Russia\n\nНомер карты " +
		"Сбер: %s\nНомер карты Тинькофф: %s\nВ комментарии укажи номер заказа [%s]\n\nПосле оплаты нажми кнопку «Оплачено»\n"

	cartWarn = "Вся информация, которую ты указываешь, для сборки в корзине 🧺 должна быть актуальной, " +
		"если она составляет более 48ч ⌚️и является неактуальной  – заказ не будет принят и деньги возвратятся в полном " +
		"объеме на карту плательщика 💴\n\n%s - Твоя заявка готова!\n\n"

	orderStartTemplate = "Номер заказа: [%s]\n\nДанные " +
		"получателя:\nФИО: %s\nНомер телефона: %s\nАдрес доставки: %s\n\nТоваров в корзине: %d\n\n"

	orderEndTemplate = "Итоговая стоимость составляет %d ₽\n"

	sendingRequsitiesTemplate = "\nВысылаю реквизиты для оплаты 🧾"

	successfulPaymentTemplate = "%s, твой заказ %s сейчас на подтверждении у админа. Он напишет тебе в личные сообщения " +
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
	return askForFIOTemplate
}

func InvalidFIO() string {
	return invalidFIOInputTemplate
}

func AskForPhoneNumber() string {
	return askForPhoneNumberTemplate
}

func AskForDeliveryAddress() string {
	return askForDeliveryAddressTemplate
}

func Requisites(shortOrderID string, r domain.Requisites) string {
	return fmt.Sprintf(requisitesTemplate, shortOrderID, r.SberID, r.TinkoffID, shortOrderID)
}

func RenderOrderAfterPayment(order domain.HouseholdOrder) string {
	start := cartWarn + _orderStart(order)
	for i, pos := range order.Cart {
		start += _cartPositionTemplate(cartPositionArgs{
			n:         i + 1,
			price:     pos.Price,
			priceGlob: pos.PriceGlob,
			name:      pos.Name,
			isbn:      pos.ISBN,
		})
	}
	return start + _orderEnd(order.AmountRUB) + sendingRequsitiesTemplate
}

func RenderMyOrders(name string, orders []domain.HouseholdOrder) string {
	start := fmt.Sprintf(myOrdersStart, name)
	for _, o := range orders {
		start += _myOrderBody(o)
		for i, pos := range o.Cart {
			start += _cartPositionTemplate(cartPositionArgs{
				n:         i + 1,
				price:     pos.Price,
				priceGlob: pos.PriceGlob,
				name:      pos.Name,
				isbn:      pos.ISBN,
			})
		}
		start += _orderEnd(o.AmountRUB) + myOrdersSeparator
	}
	return start
}

func SuccessfulPayment(fullName, orderShortID string) string {
	return fmt.Sprintf(successfulPaymentTemplate, fullName, orderShortID)
}

func _orderStart(order domain.HouseholdOrder) string {
	return fmt.Sprintf(
		orderStartTemplate,
		order.ShortID,
		*order.Customer.FullName,
		*order.Customer.PhoneNumber,
		order.DeliveryAddress,
		len(order.Cart),
	)
}

func _orderEnd(totalRub uint32) string {
	return fmt.Sprintf(orderEndTemplate, totalRub)
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
