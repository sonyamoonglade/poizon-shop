package templates

import (
	"fmt"

	"domain"
)

const (
	productTemplate         = "Название: %s\n\nЦена: %d ₽\nЦена по рынку: %d ₽\nАртикул: *%s*\n\nОписание:\n - %s"
	productTemplateDiscount = "Название: %s\n\nЦена: %d ₽\nЦена со скидкой: %d ₽\nЦена по рынку: %d ₽\nАртикул: *%s*\n\nОписание:\n - %s"
	positionAddedTemplate   = "Позиция %s успешно добавлена"
)

func HouseholdProductCaption(hp domain.HouseholdProduct) string {
	return fmt.Sprintf(productTemplate,
		hp.Name,
		hp.Price,
		hp.PriceGlob,
		hp.ISBN,
		hp.Settings,
	)
}

func HouseholdProductCaptionWithDiscount(hp domain.HouseholdProduct, discount uint32) string {
	return fmt.Sprintf(productTemplateDiscount,
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
