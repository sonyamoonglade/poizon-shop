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

	askForButtonColorTemplate = "Выбери цвет кнопки\n(влияет на условия доставки 🚚 и цену 🥬 в дальнейшем)"

	askForSizeTemplate = "Шаг 1. Выбери размер 📏\nНапример: L или 54\nЕсли товар безразмерный, то отправь #"

	askForPriceTemplate = "Отправь стоимость товара в юанях ¥\n(указана на выбранной кнопке) 💴"

	askForCategoryTemplate = "Выбери категорию товара (влияет на итоговую стоимость) 💴\n\n" +
		"В категорию «легкой одежды»относится вся обувь, кроме зимней и одежда, кроме курток 👟🧢\n\n" +
		"В категорию «тяжелая одежда»относятся все куртки и зимняя обувь 🧥🥾"

	askForLinkTemplate = "Отправь ссылку на выбранный товар (строго по инструкции) 📝"

	askForCalculatorInputTemplate = "Отправь стоимость товара в ¥, я посчитаю это в ₽  🇨🇳🇷🇺\n\n" +
		"Стоимость указана с учетом доставки товара из Китая до Москвы, доставка в другие города и " +
		"районы России просчитывается и оплачивается отдельно в ТК СДЕК 🚚"

	editPositionTemplate = "Выбери номер позиции, чтобы удалить её 🙅‍♂️\n\nПо клику на " +
		"кнопку позиция изчезнет из твоей корзины!"

	newPositionWarnTemplate = "Новый добавленный товар будет соответствовать типу доставки первоначально добавленного товара в корзине 🦧"

	deliveryOnlyToMoscowTemplate = "Стоимость указана с учетом доставки товара из Китая до Москвы, доставка в другие " +
		"города и районы России просчитывается и оплачивается отдельно в ТК СДЕК 🚚"

	requisitesTemplate = "Счет для оплаты заказа: [%s]\n\nТимофеев Вадим Денисович \\uD83D\\uDC81\\u200D♂️ @xKK_Russia\n\nНомер карты Сбер: %s\nНомер карты Тинькофф: %s\nВ комментарии укажи номер заказа [%s]\n\nПосле оплаты нажми кнопку «Оплачено»\n"
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
