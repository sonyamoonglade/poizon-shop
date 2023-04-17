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

	orderStartTemplate = "Вся информация, которую ты указываешь, для сборки в корзине 🧺 должна быть актуальной, " +
		"если она составляет более 48ч ⌚️и является неактуальной  – заказ не будет принят и деньги возвратятся в полном " +
		"объеме на карту плательщика 💴\n\n%s - Твоя заявка готова!\n\nНомер заказа: [%s]\n\nДанные " +
		"получателя\nФИО: %s\nНомер телефона: %s\nАдрес доставки: %s\n\nТоваров в корзине: %d\n\n"

	orderEndTemplate = "Итоговая стоимость составляет %d ₽\n\nВысылаю реквизиты для оплаты 🧾"

	successfulPaymentTemplate = "%s, твой заказ %s сейчас на подтверждении у админа. Он напишет тебе в личные сообщения " +
		"и подтвердит статус покупки.\n\n‼️Никому кроме бота деньги отправлять не нужно‼️Даже админу‼️\n\nТолько " +
		"админ проверяет поступление денег и обозначает статус покупки ✅"
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

func AskForOrderType() string {
	return askForOrderTypeTemplate
}

func RenderOrder(order domain.HouseholdOrder) string {
	start := _orderStart(order)
	for _, pos := range order.Cart {
		start += _cartPositionTemplate(cartPositionArgs{
			n:           len(order.Cart),
			priceRub:    pos.Price,
			productName: pos.Name,
		})
	}
	return start + _orderEnd(order.AmountRUB)
}

func SuccessfulPayment(fullName, orderShortID string) string {
	return fmt.Sprintf(successfulPaymentTemplate, fullName, orderShortID)
}

func _orderStart(order domain.HouseholdOrder) string {
	return fmt.Sprintf(
		orderStartTemplate,
		*order.Customer.FullName,
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

func _formatIsExpress(isExpress bool) string {
	if isExpress {
		return domain.ExpressStr
	}
	return domain.NormalStr
}
