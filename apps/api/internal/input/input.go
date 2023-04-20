package input

import (
	"domain"
	"dto"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddCommentToOrderInput struct {
	Comment string `json:"comment"`
}

func (a AddCommentToOrderInput) ToDTO(orderID primitive.ObjectID) dto.AddCommentDTO {
	return dto.AddCommentDTO{
		OrderID: orderID,
		Comment: a.Comment,
	}
}

type ChangeOrderStatusInput struct {
	NewStatus int `json:"newStatus"`
}

func (c ChangeOrderStatusInput) ToDTO(orderID primitive.ObjectID) dto.ChangeOrderStatusDTO {
	return dto.ChangeOrderStatusDTO{
		OrderID:   orderID,
		NewStatus: domain.Status(c.NewStatus),
	}
}

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
