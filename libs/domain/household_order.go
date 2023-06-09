package domain

import (
	"time"

	"functools"
	fn "github.com/sonyamoonglade/go_func"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HouseholdOrder struct {
	OrderID          primitive.ObjectID `json:"orderId,omitempty" bson:"_id,omitempty"`
	ShortID          string             `json:"shortId" bson:"shortId"`
	Customer         HouseholdCustomer  `json:"customer" bson:"customer"`
	Cart             HouseholdCart      `json:"cart" bson:"cart"`
	AmountRUB        uint32             `json:"amountRub" bson:"amountRub"`
	DiscountedAmount uint32             `json:"discountedAmount" bson:"discountedAmount"`
	DeliveryAddress  string             `json:"deliveryAddress" bson:"deliveryAddress"`
	Comment          *string            `json:"comment" bson:"comment"`
	InStock          bool               `json:"inStock" bson:"inStock"`
	Status           Status             `json:"status" bson:"status"`
	Source           Source             `json:"source" bson:"source"`
	CreatedAt        time.Time          `json:"createdAt" bson:"createdAt"`
	IsPaid           bool               `json:"isPaid" bson:"isPaid"`
	IsApproved       bool               `json:"isApproved" bson:"isApproved"`
}

func NewHouseholdOrder(
	customer HouseholdCustomer,
	deliveryAddress string,
	shortID string,
	inStock bool,
) HouseholdOrder {

	amountRub := functools.Reduce(func(acc uint32, cartItem HouseholdProduct) uint32 {
		acc += cartItem.Price
		return acc
	}, customer.Cart, 0)

	return HouseholdOrder{
		Customer:         customer,
		ShortID:          shortID,
		Cart:             customer.Cart,
		AmountRUB:        amountRub,
		DiscountedAmount: amountRub,
		DeliveryAddress:  deliveryAddress,
		Status:           StatusNotApproved,
		Source:           SourceHousehold,
		InStock:          inStock,
		CreatedAt:        time.Now().UTC(),
	}
}

func (h *HouseholdOrder) UseDiscount(discount uint32) HouseholdOrder {
	h.DiscountedAmount = fn.
		Reduce(func(acc uint32, el HouseholdProduct, _ int) uint32 {
			return acc + el.Price - discount
		}, h.Cart.Slice(), 0)
	return *h
}

const defaultComment = "комментарий админа отсутствует"

func (h HouseholdOrder) GetComment() string {
	if h.Comment != nil {
		return *h.Comment
	}
	return defaultComment
}
