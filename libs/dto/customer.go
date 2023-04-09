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
