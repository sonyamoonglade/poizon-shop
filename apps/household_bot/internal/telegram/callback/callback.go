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
	// Must start with this to save up symbols. See callback.Inject
	SelectCategory
	SelectSubcategory
	SelectProduct
	FromProductCardToProducts
	// --
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

func Inject(cb Callback, values ...string) string {
	if len(values) == 0 {
		return cb.string()
	}
	out := cb.string() + ";" + strings.Join(transliterators.Encode(values), ";")
	if len(out) > 64 {
		return NoOpCallback.string()
	}
	return out
}

func ParseButtonData(data string) (Callback, []string, error) {
	splitBySep := strings.Split(data, ";")
	// Nothing has been encoded
	if len(splitBySep) == 1 {
		cb, err := parseCallback(data)
		if err != nil {
			return -1, nil, fmt.Errorf("callback parse: %w", err)
		}
		return cb, nil, nil
	}
	values := transliterators.Decode(splitBySep)
	if len(values) == 0 {
		return -1, nil, fmt.Errorf("invalid callback data")
	}
	cb, err := parseCallback(values[0])
	if err != nil {
		return -1, nil, fmt.Errorf("callback parse: %w", err)
	}
	return cb, values[1:], nil
}

func parseCallback(strCallback string) (Callback, error) {
	v, err := strconv.Atoi(strCallback)
	if err != nil {
		return 0, err
	}
	return Callback(v), nil
}
