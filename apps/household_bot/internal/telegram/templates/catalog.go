package templates

import (
	"fmt"

	"domain"
)

const (
	productTemplateOrdered         = "Название: %s\n\nЦена: %d ₽\nЦена по рынку: %d ₽\n\nАртикул: *%s*\n\nОписание:\n - %s"
	productTemplateInStock         = "Название: %s\n\nЦена: %d ₽\nЦена по рынку: %d ₽\n\nНаличие: %s\nАртикул: *%s*\n\nОписание:\n - %s"
	productTemplateDiscountOrdered = "Название: %s\n\nЦена: %d ₽\nЦена со скидкой: %d ₽\nЦена по рынку: %d ₽\n\nАртикул: *%s*\n\nОписание:\n - %s"
	productTemplateDiscountInStock = "Название: %s\n\nЦена: %d ₽\nЦена со скидкой: %d ₽\nЦена по рынку: %d ₽\n\nНаличие: %s\nАртикул: *%s*\n\nОписание:\n - %s"
	positionAddedTemplate          = "Позиция %s успешно добавлена"
)

func HouseholdProductCaption(hp domain.HouseholdProduct, inStock bool) string {
	if inStock {
		return fmt.Sprintf(productTemplateInStock,
			hp.Name,
			hp.Price,
			hp.PriceGlob,
			hp.GetAvailableInStr(),
			hp.ISBN,
			hp.Settings,
		)
	}
	return fmt.Sprintf(productTemplateOrdered,
		hp.Name,
		hp.Price,
		hp.PriceGlob,
		hp.ISBN,
		hp.Settings,
	)
}

func HouseholdProductCaptionWithDiscount(hp domain.HouseholdProduct, discount uint32, inStock bool) string {
	if inStock {
		return fmt.Sprintf(productTemplateDiscountInStock,
			hp.Name,
			hp.Price,
			hp.Price-discount,
			hp.PriceGlob,
			hp.GetAvailableInStr(),
			hp.ISBN,
			hp.Settings,
		)
	}
	return fmt.Sprintf(productTemplateDiscountOrdered,
		hp.Name,
		hp.Price,
		hp.Price-discount,
		hp.PriceGlob,
		hp.ISBN,
		hp.Settings,
	)
}

func PositionAdded(name string) string {
	return fmt.Sprintf(positionAddedTemplate, name)
}
