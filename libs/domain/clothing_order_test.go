package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func stringPtr(s string) *string {
	return &s
}

func TestNewOrder(t *testing.T) {
	customer1 := Customer{
		TelegramID:  123,
		Username:    stringPtr("john"),
		FullName:    stringPtr("John Doe"),
		PhoneNumber: stringPtr("123456789"),
		TgState:     State{V: 1},
		Cart: cart.Cart{
			cart.Position{
				ShopLink:  "example.com",
				PriceRUB:  100,
				PriceYUAN: 10,
				Button:    cart.Button95,
				Size:      "xl",
			},
			cart.Position{
				ShopLink:  "example.com",
				PriceRUB:  200,
				PriceYUAN: 20,
				Button:    cart.ButtonGrey,
				Size:      "L",
			},
		},
	}

	tests := []struct {
		description     string
		customer        Customer
		deliveryAddress string
		expectedOrder   Order
	}{
		{
			description:     "test with empty deliveryAddress",
			customer:        customer1,
			deliveryAddress: "",
			expectedOrder: Order{
				Customer:        customer1,
				Cart:            customer1.Cart,
				AmountRUB:       300,
				AmountYUAN:      30,
				DeliveryAddress: "",
				IsPaid:          false,
				IsApproved:      false,
				Status:          StatusNotApproved,
			},
		},
		{
			description:     "test with non-empty deliveryAddress",
			customer:        customer1,
			deliveryAddress: "123 Main St., Anytown, USA",
			expectedOrder: Order{
				Customer:        customer1,
				Cart:            customer1.Cart,
				AmountRUB:       300,
				AmountYUAN:      30,
				DeliveryAddress: "123 Main St., Anytown, USA",
				IsPaid:          false,
				IsApproved:      false,
				Status:          StatusNotApproved,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			order := NewOrder(test.customer, test.deliveryAddress, false, "abcd")
			require.Equal(t, test.expectedOrder.Status, order.Status)
			require.Equal(t, test.expectedOrder.AmountRUB, order.AmountRUB)
			require.Equal(t, test.expectedOrder.AmountYUAN, order.AmountYUAN)
			require.Equal(t, test.expectedOrder.DeliveryAddress, order.DeliveryAddress)
			require.Equal(t, test.expectedOrder.Cart, order.Cart)
			require.False(t, order.IsPaid)
			require.False(t, order.IsApproved)
		})
	}
}
