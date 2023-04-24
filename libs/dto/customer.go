package dto

import (
	"domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UpdateClothingCustomerDTO struct {
	Username       *string
	FullName       *string
	PhoneNumber    *string
	CatalogOffset  *uint
	Meta           *domain.Meta
	CalculatorMeta *domain.CalculatorMeta
	Cart           *domain.ClothingCart
	LastPosition   *domain.ClothingPosition
	State          *domain.State
	PromocodeID    *primitive.ObjectID
}

type UpdateHouseholdCustomerDTO struct {
	Username    *string
	FullName    *string
	PhoneNumber *string
	Cart        *domain.HouseholdCart
	State       *domain.State
	PromocodeID *primitive.ObjectID
}
