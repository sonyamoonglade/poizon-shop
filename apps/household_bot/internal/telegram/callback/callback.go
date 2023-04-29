package callback

import (
	"fmt"
	"strconv"
	"strings"

	"utils/transliterators"
)

type Callback int

const (
	NoOpCallback Callback = iota
	Menu
	Catalog
	MyOrders
	GetProductByISBN
	MyCart
	Promocode
	Faq
	GetFaqAnswer
	CTypeOrder
	CTypeInStock
	SelectCategory
	SelectSubcategory
	SelectProduct
	FromProductCardToProducts
	AddToCart
	AddToCartByISBN
	EditCart
	DeletePositionFromCart
	MakeOrder
	SelectOrderType
	AcceptPayment
)

func (c Callback) string() string {
	return strconv.Itoa(int(c))
}

const (
	dataPrefix        = "d"
	rawCallbackPrefix = "c"
)

func Inject(cb Callback, values ...string) string {
	if len(values) == 0 {
		return rawCallbackPrefix + cb.string()
	}
	out := dataPrefix + cb.string() + ";" + strings.Join(transliterators.Encode(values), ";")
	if len(out) > 64 {
		return rawCallbackPrefix + NoOpCallback.string()
	}
	return out
}

func ParseButtonData(data string) (Callback, []string, error) {
	if data[0] == rawCallbackPrefix[0] {
		cb, err := parse(data[1:])
		if err != nil {
			return -1, nil, fmt.Errorf("callback parse: %w", err)
		}
		return cb, nil, nil
	}
	values := transliterators.Decode(strings.Split(data[1:], ";"))
	if len(values) == 0 {
		return -1, nil, fmt.Errorf("invalid callback data")
	}
	// Remove prefix
	cb, err := parse(values[0])
	if err != nil {
		return -1, nil, fmt.Errorf("callback parse: %w", err)
	}
	return cb, values[1:], nil
}

func parse(strCallback string) (Callback, error) {
	v, err := strconv.Atoi(strCallback)
	if err != nil {
		return 0, err
	}
	return Callback(v), nil
}
