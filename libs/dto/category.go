package dto

import "domain"

type UpdateCategoryDTO struct {
	Title         *string               `json:"title"`
	Subcategories *[]domain.Subcategory `json:"subcategories"`
	Rank          *uint                 `json:"rank"`
	Active        *bool                 `json:"active"`
}
