package dto

import "domain"

type UpdateCategoryDTO struct {
	Title         *string
	Subcategories *[]domain.Subcategory
	Rank          *uint
	Active        *bool
}
