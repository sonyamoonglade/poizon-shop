package dto

import (
	"domain"
)

type UpdateClothingCustomerDTO struct {
	LastPosition   *domain.ClothingPosition
	Username       *string
	FullName       *string
	Meta           *domain.Meta
	CalculatorMeta *domain.CalculatorMeta
	PhoneNumber    *string
	Cart           *domain.ClothingCart
	State          *domain.State
	CatalogOffset  *uint
}

type UpdateHouseholdCustomerDTO struct {
	Username    *string
	FullName    *string
	PhoneNumber *string
	Meta        *domain.Meta
	Cart        *domain.HouseholdCart
	State       *domain.State
}
