package domain

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrNoCategories     = errors.New("no categories")
)

type HouseholdCategory struct {
	CategoryID primitive.ObjectID `json:"categoryId" bson:"_id,omitempty"`
	Title      string             `json:"title" bson:"title"`
	Active     bool               `json:"active" bson:"active"`
	// Immutable field
	InStock bool `json:"inStock" bson:"inStock"`
	// Less - first
	Rank          uint          `json:"rank" bson:"rank"`
	Subcategories []Subcategory `json:"subcategories" bson:"subcategories"`
}

func NewHouseholdCategory(title string, inStock bool) HouseholdCategory {
	return HouseholdCategory{
		Title:   title,
		InStock: inStock,
	}
}

func (c *HouseholdCategory) SetRank(r uint) *HouseholdCategory {
	c.Rank = r
	return c
}

const (
	in  = "В наличии"
	ord = "Под заказ"
)

func InStockToString(inStock bool) string {
	if inStock {
		return in
	}
	return ord
}
