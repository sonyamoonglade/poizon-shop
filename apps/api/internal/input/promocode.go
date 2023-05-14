package input

import (
	"domain"
	"dto"
)

type NewPromocodeInput struct {
	Description string                                 `json:"description"`
	Discounts   map[string] /* source string */ uint32 `json:"discounts"`
	ShortID     string                                 `json:"shortId"`
}

func (n NewPromocodeInput) ToDomainPromocode() domain.Promocode {
	discounts := make(domain.DiscountMap)
	for src, discount := range n.Discounts {
		discounts[domain.SourceFromString(src)] = discount
	}
	return domain.NewPromocode(n.Description, discounts, n.ShortID)
}

type UpdatePromocodeInput struct {
	Discounts map[string]uint32 `json:"discounts"`
}

func (u UpdatePromocodeInput) ToDTO() dto.UpdatePromocodeDTO {
	discounts := make(domain.DiscountMap)
	for src, discount := range u.Discounts {
		discounts[domain.SourceFromString(src)] = discount
	}
	return dto.UpdatePromocodeDTO{
		Discounts: &discounts,
	}
}
