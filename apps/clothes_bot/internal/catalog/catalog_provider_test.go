package catalog

import (
	"sync"
	"testing"

	"domain"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCatalogProvider(t *testing.T) {
	items := []domain.ClothingProduct{
		{
			ItemID:         primitive.NewObjectID(),
			ImageURLs:      []string{"https://example.com/image.jpg"},
			Title:          "Example Item",
			Rank:           1,
			AvailableSizes: []string{"XS", "S"},
		},
		{
			ItemID:    primitive.NewObjectID(),
			ImageURLs: []string{"https://example.com/image.jpg"},
			Title:     "Example Item 2",
			Rank:      2,

			AvailableSizes: []string{"XS", "S"},
		},
		{
			ItemID:         primitive.NewObjectID(),
			ImageURLs:      []string{"https://example.com/image.jpg"},
			Title:          "Example Item 3",
			Rank:           3,
			AvailableSizes: []string{"XS", "S"},
		},
	}

	cp := CatalogProvider{
		mu:    &sync.RWMutex{},
		items: items,
	}

	t.Run("test has next", func(t *testing.T) {
		tests := []struct {
			description string
			offset      uint
			expected    bool
		}{
			{"offset = 0, should return true because elem [1] exists", 0, true},
			{"offset = 1, should return true because elem[2] exists", 1, true},
			{"offset = 2, should return false becase elem[3] no exists", 2, false},
		}

		for _, test := range tests {
			t.Run(test.description, func(t *testing.T) {
				actual := cp.HasNext(test.offset)
				require.Equal(t, test.expected, actual)
			})
		}
	})

	t.Run("test has prev", func(t *testing.T) {
		tests := []struct {
			description string
			offset      uint
			expected    bool
		}{
			{"offset = 0, should return false because no elem under elem[-1]", 0, false},
			{"offset = 1, should return true because elem[0] exists", 1, true},
			{"offset = 2, should return true because elem[1] exists", 2, true},
			{"offset = 3, should return true because elem[2] exists", 3, true},
		}

		for _, test := range tests {
			actual := cp.HasPrev(test.offset)
			t.Run(test.description, func(t *testing.T) {
				require.Equal(t, test.expected, actual)
			})
		}
	})

	t.Run("test load next", func(t *testing.T) {
		tests := []struct {
			description string
			offset      uint
			expected    domain.ClothingProduct
		}{
			{"offset = 0, should load 2nd elem", 0, items[1]},
			{"offset = 1, should load 3rd elem", 1, items[2]},
			{"offset = 2 should be empty", 2, domain.ClothingProduct{}},
		}
		for _, test := range tests {
			actual := cp.LoadNext(test.offset)
			t.Run(test.description, func(t *testing.T) {
				require.Equal(t, test.expected, actual)
			})
		}
	})

	t.Run("test load prev", func(t *testing.T) {
		tests := []struct {
			description string
			offset      uint
			expected    domain.ClothingProduct
		}{
			{"offset = 0, should return empty", 0, domain.ClothingProduct{}},
			{"offset = 1, should load 1st elem", 1, items[0]},
			{"offset = 2, shoud load 2nd elem", 2, items[1]},
		}
		for _, test := range tests {
			actual := cp.LoadPrev(test.offset)
			t.Run(test.description, func(t *testing.T) {
				require.EqualValues(t, test.expected, actual)
			})
		}
	})

	t.Run("load first", func(t *testing.T) {
		tests := []struct {
			description string
			expected    domain.ClothingProduct
		}{
			{"should return first", items[0]},
		}
		for _, test := range tests {
			actual := cp.LoadFirst()
			t.Run(test.description, func(t *testing.T) {
				require.EqualValues(t, test.expected, actual)
			})
		}
	})
}
