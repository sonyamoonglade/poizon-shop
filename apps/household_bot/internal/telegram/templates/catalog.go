package templates

import (
	"fmt"

	"domain"
)

const (
	productTemplate       = "nazv: %s\nisbn: %s\nprice: %d\npriceGlob: %d\n\nsettings: %s"
	positionAddedTemplate = "Позиция %s успешно добавлена"
)

func HouseholdProductCaption(hp domain.HouseholdProduct) string {
	return fmt.Sprintf(productTemplate,
		hp.Name,
		hp.ISBN,
		hp.Price,
		hp.PriceGlob,
		hp.Settings,
	)
}

func PositionAdded(name string) string {
	return fmt.Sprintf(positionAddedTemplate, name)
}

func formatBool(b bool) string {
	if b {
		return "Da"
	}
	return "No"
}
