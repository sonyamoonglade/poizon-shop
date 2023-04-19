package tests

import (
	"context"
	"strconv"

	"domain"
	f "github.com/brianvoe/gofakeit/v6"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *AppTestSuite) TestHandlerAddToCart() {
	var (
		require    = s.Require()
		telegramID = f.Int64()
		username   = f.Username()
		ctx        = context.Background()
	)

	s.Run("should not add because customer has added"+
		"item with category (inStock=false), but trying to add"+
		"product with category(inStock=true)", func() {

		// Product from this category is already added to cart (in stock)
		c1 := domain.NewHouseholdCategory(f.Word(), true)
		c1.CategoryID = primitive.NewObjectID()
		c1.Subcategories = append(c1.Subcategories, domain.Subcategory{
			SubcategoryID: primitive.NewObjectID(),
			Title:         f.Word(),
			Active:        true,
			Rank:          0,
			Products: []domain.HouseholdProduct{
				{
					CategoryID: c1.CategoryID,
					ImageURL:   f.ImageURL(200, 300),
					Name:       f.StreetName(),
					ISBN:       f.BuzzWord(),
					Settings:   f.LoremIpsumWord(),
					Price:      1223,
					PriceGlob:  1231,
				},
			},
		})

		// customer tries to add product from this category (not in stock)
		c2 := domain.NewHouseholdCategory(f.Word(), false)
		c2.CategoryID = primitive.NewObjectID()
		c2.Subcategories = append(c2.Subcategories, domain.Subcategory{
			SubcategoryID: primitive.NewObjectID(),
			Title:         f.Word(),
			Active:        true,
			Rank:          0,
			Products: []domain.HouseholdProduct{
				{
					CategoryID: c2.CategoryID,
					ImageURL:   f.ImageURL(200, 300),
					Name:       f.StreetName(),
					ISBN:       f.BuzzWord(),
					Settings:   f.LoremIpsumWord(),
					Price:      1223,
					PriceGlob:  1231,
				},
			},
		})

		err := s.repositories.HouseholdCategory.Save(ctx, c1)
		require.NoError(err)
		err = s.repositories.HouseholdCategory.Save(ctx, c2)
		require.NoError(err)

		customer := domain.NewHouseholdCustomer(telegramID, username)
		// Initially has this product in cart
		customer.Cart.Add(c1.Subcategories[0].Products[0])
		err = s.repositories.HouseholdCustomer.Save(ctx, customer)
		require.NoError(err)

		// c2,product from c2
		var (
			cTitle     = c2.Title
			sTitle     = c2.Subcategories[0].Title
			inStockStr = strconv.FormatBool(c2.InStock)
			pName      = c2.Subcategories[0].Products[0].Name
		)
		expectedArgs := []string{cTitle, sTitle, inStockStr, pName}

		err = s.tghandler.AddToCart(ctx, telegramID, expectedArgs)
		require.NoError(err)

		customer, err = s.repositories.HouseholdCustomer.GetByTelegramID(ctx, telegramID)
		require.NoError(err)

		// Second product not added
		require.True(len(customer.Cart) == 1)

		s.repositories.HouseholdCustomer.Delete(ctx, customer.CustomerID)
		s.repositories.HouseholdCategory.Delete(ctx, c1.CategoryID)
		s.repositories.HouseholdCategory.Delete(ctx, c2.CategoryID)
	})

}
