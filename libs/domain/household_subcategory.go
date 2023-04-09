package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Subcategory struct {
	SubcategoryID primitive.ObjectID `json:"subcategoryId" bson:"_id,omitempty"`
	Title         string             `json:"title" bson:"title"`
	Active        bool               `json:"active" bson:"active"`
	Rank          uint               `json:"rank" bson:"rank"`
	Products      []HouseholdProduct `json:"products" bson:"products"`
}
