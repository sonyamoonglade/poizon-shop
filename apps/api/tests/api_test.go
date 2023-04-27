package tests

import (
	"context"
	"encoding/json"
	"net/http"

	"domain"
	f "github.com/brianvoe/gofakeit/v6"
	"onlineshop/api/internal/input"
	"utils/testutil"
)

func (s *AppTestSuite) TestApiAddItem() {
	var (
		require = s.Require()
	)

	s.Run("should add item to catalog with rank 0 because catalog's empty", func() {
		inp := input.AddItemToCatalogInput{
			ImageURLs:       []string{f.URL(), f.URL()},
			AvailableSizes:  []string{f.StreetNumber(), f.StreetNumber()},
			AvailableInCity: []string{f.City(), f.City()},
			Title:           f.BeerName(),
			Quantity:        f.IntRange(1, 10),
			ShopLink:        f.URL(),
			PriceRUB:        234,
		}
		resp, err := s.app.Test(testutil.NewJsonRequestWithKey("abcd")(http.MethodPost, "/api/clothing/catalog/addItem", inp), -1)
		require.NoError(err)
		require.Equal(http.StatusOK, resp.StatusCode)

		var respJson []domain.ClothingProduct
		require.NoError(json.NewDecoder(resp.Body).Decode(&respJson))
		require.True(len(respJson) == 1)
		elem := respJson[0]
		require.Equal(uint(0), elem.Rank)

		s.repositories.ClothingCatalog.RemoveItem(context.Background(), elem.ItemID)
	})

	s.Run("should add item to catalog with rank 0 and then next one with 1 sequentially", func() {
		var inputs = []input.AddItemToCatalogInput{
			{
				ImageURLs:       []string{f.URL(), f.URL()},
				AvailableSizes:  []string{f.StreetNumber(), f.StreetNumber()},
				AvailableInCity: []string{f.City(), f.City()},
				Title:           f.Username(),
				Quantity:        f.IntRange(1, 10),
				ShopLink:        f.URL(),
				PriceRUB:        235,
			},
			{
				ImageURLs:       []string{f.URL(), f.URL()},
				AvailableSizes:  []string{f.StreetNumber(), f.StreetNumber()},
				AvailableInCity: []string{f.City(), f.City()},
				Title:           f.Username(),
				Quantity:        f.IntRange(1, 10),
				ShopLink:        f.URL(),
				PriceRUB:        236,
			},
			{
				ImageURLs:       []string{f.URL(), f.URL()},
				AvailableSizes:  []string{f.StreetNumber(), f.StreetNumber()},
				AvailableInCity: []string{f.City(), f.City()},
				Title:           f.Username(),
				Quantity:        f.IntRange(1, 10),
				ShopLink:        f.URL(),
				PriceRUB:        237,
			},
		}
		i1, i2, i3 := inputs[0], inputs[1], inputs[2]

		// Add first
		resp, err := s.app.Test(testutil.NewJsonRequestWithKey("abcd")(http.MethodPost, "/api/clothing/catalog/addItem", i1), -1)
		require.NoError(err)
		require.Equal(http.StatusOK, resp.StatusCode)

		// Add second
		resp, err = s.app.Test(testutil.NewJsonRequestWithKey("abcd")(http.MethodPost, "/api/clothing/catalog/addItem", i2), -1)
		require.NoError(err)
		require.Equal(http.StatusOK, resp.StatusCode)

		var respJson []domain.ClothingProduct
		// Add third
		resp, err = s.app.Test(testutil.NewJsonRequestWithKey("abcd")(http.MethodPost, "/api/clothing/catalog/addItem", i3), -1)
		require.NoError(err)
		require.Equal(http.StatusOK, resp.StatusCode)
		require.NoError(json.NewDecoder(resp.Body).Decode(&respJson))

		var elem1, elem2, elem3 domain.ClothingProduct
		for _, item := range respJson {
			if item.Title == inputs[0].Title {
				elem1 = item
			}
			if item.Title == inputs[1].Title {
				elem2 = item
			}
			if item.Title == inputs[2].Title {
				elem3 = item
			}
		}
		// Added firstly
		require.Equal(uint(0), elem1.Rank)
		// Added secondly
		require.Equal(uint(1), elem2.Rank)
		require.Equal(uint(2), elem3.Rank)
		require.Equal(3, len(respJson))

		// cleanup
		s.repositories.ClothingCatalog.RemoveItem(context.Background(), elem1.ItemID)
		s.repositories.ClothingCatalog.RemoveItem(context.Background(), elem2.ItemID)
		s.repositories.ClothingCatalog.RemoveItem(context.Background(), elem3.ItemID)
	})
}

func (s *AppTestSuite) TestDeleteItem() {
	var (
		require = s.Require()
	)

	s.Run("should add 3 catalog items, remove last one, ranks should remain the same and customer's offsets set to 0", func() {
		// Add items
		ctx := context.Background()
		for i := 0; i < 3; i++ {
			item := catalogItemFixture()
			item.Rank = uint(i)
			err := s.repositories.ClothingCatalog.AddItem(ctx, item)
			require.NoError(err)
		}

		catalog, err := s.repositories.ClothingCatalog.GetCatalog(ctx)
		require.NoError(err)
		// remove item with top rank (last)
		deleteItem := catalog[len(catalog)-1]
		resp, err := s.app.Test(testutil.NewJsonRequestWithKey("abcd")(http.MethodPost, "/api/clothing/catalog/deleteItem", input.RemoveItemFromCatalogInput{
			ItemID: deleteItem.ItemID,
		}), -1)
		require.NoError(err)
		require.Equal(http.StatusOK, resp.StatusCode)

		var respJson []domain.ClothingProduct
		require.NoError(json.NewDecoder(resp.Body).Decode(&respJson))
		// Check if deleted
		for _, newItem := range respJson {
			if newItem.ItemID == deleteItem.ItemID {
				require.FailNowf("failed", "item with id: %s has not been deleted", deleteItem.ItemID.String())
			}
		}

		require.True(respJson[0].Rank == 0)
		require.True(respJson[1].Rank == 1)

		customers, err := s.repositories.ClothingCustomer.All(ctx)
		require.NoError(err)
		for _, c := range customers {
			require.True(c.CatalogOffset == uint(0))
		}

		// cleanup
		for _, item := range respJson {
			s.repositories.ClothingCatalog.RemoveItem(ctx, item.ItemID)
		}
	})
	s.Run("remove item in the middle. Should update ranks properly", func() {
		// Add items
		ctx := context.Background()
		for i := 0; i < 100; i++ {
			item := catalogItemFixture()
			item.Rank = uint(i)
			err := s.repositories.ClothingCatalog.AddItem(ctx, item)
			require.NoError(err)
		}

		catalog, err := s.repositories.ClothingCatalog.GetCatalog(ctx)
		require.NoError(err)
		// remove item with top rank (last)
		deleteItem := catalog[50]
		resp, err := s.app.Test(testutil.NewJsonRequestWithKey("abcd")(http.MethodPost, "/api/clothing/catalog/deleteItem", input.RemoveItemFromCatalogInput{
			ItemID: deleteItem.ItemID,
		}), -1)
		require.NoError(err)
		require.Equal(http.StatusOK, resp.StatusCode)

		var respJson []domain.ClothingProduct
		require.NoError(json.NewDecoder(resp.Body).Decode(&respJson))

		// Check if deleted and valid rank
		for i, newItem := range respJson {
			if newItem.ItemID == deleteItem.ItemID {
				require.FailNowf("failed", "item with id: %s has not been deleted", deleteItem.ItemID.String())
			}
			require.Equal(uint(i), newItem.Rank)
		}

		// cleanup
		for _, item := range respJson {
			s.repositories.ClothingCatalog.RemoveItem(ctx, item.ItemID)
		}
	})
}
