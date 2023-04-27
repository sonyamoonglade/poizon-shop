package domain

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"functools"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrNoCatalog    = errors.New("catalog not found")
	ErrItemNotFound = errors.New("item not found")
)

type ClothingProduct struct {
	ItemID          primitive.ObjectID `json:"itemId,omitempty" bson:"_id,omitempty"`
	ImageURLs       []string           `json:"imageUrls" bson:"imageUrls"`
	AvailableSizes  []string           `json:"availableSizes" bson:"availableSizes"`
	AvailableInCity []string           `json:"availableInCity" bson:"availableInCity"`
	Quantity        int                `json:"quantity" bson:"quantity"`
	Title           string             `json:"title" bson:"title"`
	ShopLink        string             `json:"shopLink" bson:"shopLink"`
	Rank            uint               `json:"rank" bson:"rank"`
	PriceRUB        uint64             `json:"priceRub" bson:"priceRub"`
}

func (c *ClothingProduct) GetSizesPretty() string {
	var out string
	for i, size := range c.AvailableSizes {
		// last
		if i == len(c.AvailableSizes)-1 {
			out += fmt.Sprintf("(%s)", size)
			continue
		}
		out += fmt.Sprintf("(%s); ", size)
	}
	return out
}

func (c *ClothingProduct) GetCitiesPretty() string {
	return strings.Join(c.AvailableInCity, "; ")
}

// catalog must be sorted by rank ascending
func UpdateRanks(catalog []ClothingProduct) []ClothingProduct {
	if catalog == nil {
		return nil
	}
	// If first item's rank is 0 then down all subsequent
	if catalog[0].Rank != uint(0) && catalog[0].Rank > uint(0) {
		return functools.Map(func(item ClothingProduct, i int) ClothingProduct {
			item.Rank--
			return item
		}, catalog)
	}

	// Found gap somewhere in between (only one at a time)
	var idxGap int
	for i := 0; i < len(catalog)-1; i++ {
		curr, next := catalog[i], catalog[i+1]
		if math.Abs(float64(curr.Rank)-float64(next.Rank)) > 1 {
			idxGap = i + 1
			break
		}
	}
	// All fine
	if idxGap == 0 {
		return catalog
	}
	return functools.Map(func(item ClothingProduct, i int) ClothingProduct {
		if i >= idxGap {
			item.Rank--
			return item
		}
		return item
	}, catalog)
}
