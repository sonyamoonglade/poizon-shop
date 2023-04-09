package domain

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
)

type HouseholdCategory struct {
	CategoryID primitive.ObjectID `json:"categoryId" bson:"_id,omitempty"`
	Title      string             `json:"title" bson:"title"`
	Active     bool               `json:"active" bson:"active"`
	// Less - first
	Rank          uint          `json:"rank" bson:"rank"`
	Subcategories []Subcategory `json:"subcategories" bson:"subcategories"`
}

func NewHouseholdCategory(title string) HouseholdCategory {
	return HouseholdCategory{
		Title: title,
	}
}

func (c HouseholdCategory) SetRank(r uint) HouseholdCategory {
	c.Rank = r
	return c
}
