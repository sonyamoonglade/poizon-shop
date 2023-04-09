package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpdateRanks(t *testing.T) {

	t.Run("nil catalog", func(t *testing.T) {
		require.NotPanics(t, func() {
			var catalog []ClothingProduct
			UpdateRanks(catalog)
		})
	})
	t.Run("delete last item. No change in ranks", func(t *testing.T) {
		var catalog []ClothingProduct
		for i := 0; i < 100; i++ {
			catalog = append(catalog, ClothingProduct{
				Rank: uint(i),
			})
		}
		// delete last
		catalog = catalog[:len(catalog)-1]
		newCatalog := UpdateRanks(catalog)
		// No change in ranks
		require.EqualValues(t, catalog, newCatalog)
	})

	t.Run("delete first item. Should update ranks by -1", func(t *testing.T) {
		var catalog []ClothingProduct
		for i := 0; i < 100; i++ {
			catalog = append(catalog, ClothingProduct{
				Rank: uint(i),
			})
		}
		// delete first
		catalog = catalog[1:]
		newCatalog := UpdateRanks(catalog)
		for i, item := range newCatalog {
			require.True(t, item.Rank == uint(i))
		}
	})

	t.Run("delete item in between. Should update ranks correctly", func(t *testing.T) {
		var catalog []ClothingProduct
		for i := 0; i < 100; i++ {
			catalog = append(catalog, ClothingProduct{
				Rank: uint(i),
			})
		}
		// delete first
		catalog = append(catalog[:50], catalog[51:]...)
		t.Log(len(catalog))
		newCatalog := UpdateRanks(catalog)
		for i, item := range newCatalog {
			require.True(t, item.Rank == uint(i))
		}
	})
}
