package dto

import "domain"

type UpdatePromocodeDTO struct {
	Description *string
	Discounts   *domain.DiscountMap
}
