package input

import (
	"domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RemoveItemFromCatalogInput struct {
	ItemID primitive.ObjectID `json:"itemId"`
}

type AddItemToCatalogInput struct {
	ImageURLs       []string `json:"imageUrls"`
	AvailableSizes  []string `json:"availableSizes"`
	AvailableInCity []string `json:"availableInCity"`
	Title           string   `json:"title"`
	Quantity        int      `json:"quantity"`
	ShopLink        string   `json:"shopLink"`
	PriceRUB        uint64   `json:"priceRub"`
}

func (a AddItemToCatalogInput) ToNewClothingProduct(rank uint) domain.ClothingProduct {
	return domain.ClothingProduct{
		ImageURLs:       a.ImageURLs,
		AvailableSizes:  a.AvailableSizes,
		Title:           a.Title,
		ShopLink:        a.ShopLink,
		AvailableInCity: a.AvailableInCity,
		Quantity:        a.Quantity,
		PriceRUB:        a.PriceRUB,
		Rank:            rank,
	}
}
