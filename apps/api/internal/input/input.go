package input

import (
	"domain"
	"dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddCommentToOrderInput struct {
	OrderID primitive.ObjectID `json:"orderId"`
	Comment string             `json:"comment"`
}

func (a AddCommentToOrderInput) ToDTO() dto.AddCommentDTO {
	return dto.AddCommentDTO{
		OrderID: a.OrderID,
		Comment: a.Comment,
	}
}

type ChangeOrderStatusInput struct {
	OrderID   primitive.ObjectID `json:"orderId"`
	NewStatus int                `json:"newStatus"`
}

func (c ChangeOrderStatusInput) ToDTO() dto.ChangeOrderStatusDTO {
	return dto.ChangeOrderStatusDTO{
		OrderID:   c.OrderID,
		NewStatus: domain.Status(c.NewStatus),
	}
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

type RemoveItemFromCatalogInput struct {
	ItemID primitive.ObjectID `json:"itemId"`
}

type RankUpInput struct {
	ItemID primitive.ObjectID `json:"itemId"`
}

type RankDownInput struct {
	ItemID primitive.ObjectID `json:"itemId"`
}
