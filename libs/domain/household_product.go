package domain

import (
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrProductsNotFound = errors.New("no products")
	ErrProductNotFound  = errors.New("product not found")
)

type HouseholdProduct struct {
	ProductID   primitive.ObjectID `json:"productId" bson:"_id,omitempty"`
	CategoryID  primitive.ObjectID `json:"categoryId" bson:"categoryId"`
	ImageURL    string             `json:"imageUrl" bson:"imageUrl"`
	Name        string             `json:"name" bson:"name"`
	ISBN        string             `json:"isbn" bson:"isbn"`
	Settings    string             `json:"settings" bson:"settings"`
	AvailableIn *[]string          `json:"availableIn,omitempty" bson:"availableIn,omitempty"`
	Price       uint32             `json:"price" bson:"price"`
	PriceGlob   uint32             `json:"priceGlob" bson:"priceGlob"`
}

func (hp HouseholdProduct) GetAvailableInStr() string {
	if hp.AvailableIn == nil {
		return ""
	}
	return strings.Join(*hp.AvailableIn, ";")
}

func (hp HouseholdProduct) HasImage() bool {
	return hp.ImageURL != ""
}
