package templates

import (
	"domain"
	"fmt"
)

const (
	productTemplate = "nazv: %s\nisbn: %s\nprice: %d\npriceGlob: %d\n\nsettings: %s\nnalichee: %s\n"
)

func HouseholdProductCaption(hp domain.HouseholdProduct) string {
	return fmt.Sprintf(productTemplate,
		hp.Name,
		hp.ISBN,
		hp.Price,
		hp.PriceGlob,
		hp.Settings,
		formatBool(hp.AvailableInStock),
	)
}

func formatBool(b bool) string {
	if b {
		return "Da"
	}
	return "No"
}
