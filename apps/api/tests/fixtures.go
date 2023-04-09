package tests

import (
	"domain"
	f "github.com/brianvoe/gofakeit/v6"
)

func catalogItemFixture() domain.ClothingProduct {
	return domain.ClothingProduct{
		ImageURLs:       []string{f.BeerName()},
		AvailableSizes:  []string{f.City(), f.City()},
		AvailableInCity: []string{f.City(), f.City()},
		Quantity:        f.IntRange(1, 10),
		Title:           f.Word(),
		ShopLink:        f.URL(),
		PriceRUB:        uint64(f.IntRange(1, 15)),
	}
}
