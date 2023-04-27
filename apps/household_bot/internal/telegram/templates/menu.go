package templates

import "fmt"

const (
	startGreeting   = "Привет, %s!\nРад видеть тебя в боте техники хКК 👋🏻 \n\nЖми на кнопку меню 👇🏻"
	askForPromocode = "Введи промокод: "
	promoWarn       = "Осторожно! Промокод можно ввести только 1 раз!"
	promoUseSuccess = "Промокод %s успешно применен!\nСумма скидки на все товары составляет %d ₽"
)

func AskForPromocode() string {
	return askForPromocode
}

func PromocodeWarning() string {
	return promoWarn
}

func PromocodeUseSuccess(shortID string, discount uint32) string {
	return fmt.Sprintf(promoUseSuccess, shortID, discount)
}

func StartGreeting(username string) string {
	return fmt.Sprintf(startGreeting, username)
}
