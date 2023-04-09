package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type HouseholdProduct struct {
	ProductID        primitive.ObjectID `json:"productId" bson:"_id,omitempty"`
	ImageURL         string             `json:"imageUrl" bson:"imageUrl"`
	ISBN             string             `json:"isbn" bson:"isbn"`
	Settings         string             `json:"settings" bson:"settings"`
	Price            uint32             `json:"price" bson:"price"`
	PriceGlob        uint32             `json:"priceGlob" bson:"priceGlob"`
	AvailableInStock bool               `json:"availableInStock" bson:"availableInStock"`
}
