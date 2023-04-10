package domain

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrProductsNotFound = errors.New("no products")
)

type HouseholdProduct struct {
	ProductID        primitive.ObjectID `json:"productId" bson:"_id,omitempty"`
	ImageURL         string             `json:"imageUrl" bson:"imageUrl"`
	Name             string             `json:"name" bson:"name"`
	ISBN             string             `json:"isbn" bson:"isbn"`
	Settings         string             `json:"settings" bson:"settings"`
	Price            uint32             `json:"price" bson:"price"`
	PriceGlob        uint32             `json:"priceGlob" bson:"priceGlob"`
	AvailableInStock bool               `json:"availableInStock" bson:"availableInStock"`
}
