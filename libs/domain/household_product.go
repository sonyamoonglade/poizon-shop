package domain

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrProductsNotFound = errors.New("no products")
	ErrProductNotFound  = errors.New("product not found")
)

type HouseholdProduct struct {
	ProductID  primitive.ObjectID `json:"productId" bson:"_id,omitempty"`
	CategoryID primitive.ObjectID `json:"categoryId" bson:"categoryId"`
	ImageURL   string             `json:"imageUrl" bson:"imageUrl"`
	Name       string             `json:"name" bson:"name"`
	ISBN       string             `json:"isbn" bson:"isbn"`
	Settings   string             `json:"settings" bson:"settings"`
	Price      uint32             `json:"price" bson:"price"`
	PriceGlob  uint32             `json:"priceGlob" bson:"priceGlob"`
}
