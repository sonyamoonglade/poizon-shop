package domain

import (
	"functools"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HouseholdOrder struct {
	OrderID         primitive.ObjectID `json:"orderId,omitempty" bson:"_id,omitempty"`
	ShortID         string             `json:"shortId" bson:"shortId"`
	Customer        HouseholdCustomer  `json:"customer" bson:"customer"`
	Cart            HouseholdCart      `json:"cart" bson:"cart"`
	AmountRUB       uint32             `json:"amountRub" bson:"amountRub"`
	DeliveryAddress string             `json:"deliveryAddress" bson:"deliveryAddress"`
	Comment         *string            `json:"comment" bson:"comment"`
	Status          Status             `json:"status" bson:"status"`
	Source          Source             `json:"source" bson:"source"`
	IsPaid          bool               `json:"isPaid" bson:"isPaid"`
	IsApproved      bool               `json:"isApproved" bson:"isApproved"`
	IsExpress       bool               `json:"isExpress" bson:"isExpress"`
}

func NewHouseholdOrder(customer HouseholdCustomer,
	deliveryAddress string,
	isExpress bool,
	shortID string) HouseholdOrder {

	amountRub := functools.Reduce(func(acc uint32, cartItem HouseholdProduct) uint32 {
		acc += cartItem.Price
		return acc
	}, customer.Cart, 0)

	return HouseholdOrder{
		Customer:        customer,
		ShortID:         shortID,
		Cart:            customer.Cart,
		AmountRUB:       amountRub,
		DeliveryAddress: deliveryAddress,
		IsExpress:       isExpress,
		Status:          StatusNotApproved,
		Source:          SourceHousehold,
	}
}
