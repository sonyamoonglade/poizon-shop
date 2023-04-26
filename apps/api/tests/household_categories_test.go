package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"domain"
	"utils/testutil"

	f "github.com/brianvoe/gofakeit/v6"
)

func (s *AppTestSuite) TestHouseholdNew() {
	var (
		require = s.Require()
	)

	s.Run("should add new category", func() {
		title := f.StreetName()
		url := fmt.Sprintf("/api/household/categories/new?title=%s", title)
		req := testutil.NewJsonRequestWithKey("abcd")(http.MethodPost, url, nil)
		res, err := s.app.Test(req, -1)
		require.NoError(err)
		require.Equal(http.StatusCreated, res.StatusCode)
		ctx := context.Background()
		categories, err := s.repositories.HouseholdCategory.GetAll(ctx)
		require.NoError(err)
		require.True(len(categories) == 1)
		c := categories[0]
		require.Equal(uint(0), c.Rank)
		require.Equal(title, c.Title)
		require.True(len(c.Subcategories) == 0)
		require.Equal(false, c.Active)

		// cleanup
		s.repositories.HouseholdCategory.Delete(ctx, c.CategoryID)
	})

	s.Run("should return all categories", func() {
		ctx := context.Background()
		for i := 0; i < 10; i++ {
			require.NoError(s.services.HouseholdCategory.New(ctx, f.Word(), true))
		}
		url := "/api/household/categories/all"
		req := testutil.NewJsonRequestWithKey("abcd")(http.MethodGet, url, nil)
		res, err := s.app.Test(req, -1)
		require.NoError(err)
		var categories []domain.HouseholdCategory
		require.NoError(json.NewDecoder(res.Body).Decode(&categories))
		require.True(len(categories) == 10)
		for i, c := range categories {
			require.Equal(uint(i), c.Rank)
		}
	})
}
