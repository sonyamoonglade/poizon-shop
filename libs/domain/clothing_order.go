package domain

import (
	"time"

	fn "github.com/sonyamoonglade/go_func"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClothingOrder struct {
	OrderID          primitive.ObjectID `json:"orderId,omitempty" bson:"_id,omitempty"`
	ShortID          string             `json:"shortId" bson:"shortId"`
	Customer         ClothingCustomer   `json:"customer" bson:"customer"`
	Cart             ClothingCart       `json:"cart" bson:"cart"`
	AmountRUB        uint64             `json:"amountRub" bson:"amountRub"`
	AmountYUAN       uint64             `json:"amountYuan" bson:"amountYuan"`
	DiscountedAmount uint64             `json:"discountedAmount" bson:"discountedAmount"`
	DeliveryAddress  string             `json:"deliveryAddress" bson:"deliveryAddress"`
	Status           Status             `json:"status" bson:"status"`
	Source           Source             `json:"source" bson:"source"`
	Comment          *string            `json:"comment" bson:"comment"`
	CreatedAt        time.Time          `json:"createdAt" bson:"createdAt"`
	IsPaid           bool               `json:"isPaid" bson:"isPaid"`
	IsApproved       bool               `json:"isApproved" bson:"isApproved"`
	IsExpress        bool               `json:"isExpress" bson:"isExpress"`
}

func NewClothingOrder(
	customer ClothingCustomer,
	deliveryAddress string,
	isExpress bool,
	shortID string,
) ClothingOrder {
	type total struct {
		rub, yuan uint64
	}
	var totals total

	totals = fn.Reduce(func(t total, cartItem ClothingPosition, _ int) total {
		t.yuan += cartItem.PriceYUAN
		t.rub += cartItem.PriceRUB
		return t
	}, customer.Cart.Slice(), total{})

	return ClothingOrder{
		Customer:         customer,
		ShortID:          shortID,
		Cart:             customer.Cart,
		AmountRUB:        totals.rub,
		AmountYUAN:       totals.yuan,
		DiscountedAmount: totals.rub,
		DeliveryAddress:  deliveryAddress,
		IsExpress:        isExpress,
		CreatedAt:        time.Now().UTC(),
		Status:           StatusNotApproved,
		Source:           SourceClothing,
	}
}

func (c *ClothingOrder) UseDiscount(discount uint32) *ClothingOrder {
	c.DiscountedAmount = fn.
		Reduce(func(acc uint64, pos ClothingPosition, _ int) uint64 {
			return acc + pos.PriceRUB - uint64(discount)
		}, c.Cart.Slice(), 0)
	return c
}

func (c ClothingOrder) GetComment() string {
	if c.Comment != nil {
		return *c.Comment
	}
	return defaultComment
}
