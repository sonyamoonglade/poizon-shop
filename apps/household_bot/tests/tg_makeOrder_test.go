package tests

import (
	"context"
	"testing"

	"domain"
	"dto"
	f "github.com/brianvoe/gofakeit/v6"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *AppTestSuite) TestHandlerMakeOrder() {
	var (
		require = s.Require()
		ctx     = context.Background()
	)

	s.Run("should not create order because not all products exist at the moment", func() {
		var (
			username   = f.Username()
			telegramID = f.Int64()
		)
		customer := domain.NewHouseholdCustomer(telegramID, username)
		c1 := domain.NewHouseholdCategory(f.Word(), true)
		c1.CategoryID = primitive.NewObjectID()
		c1.Subcategories = append(c1.Subcategories, domain.Subcategory{
			SubcategoryID: primitive.NewObjectID(),
			Title:         f.Word(),
			Active:        true,
			Rank:          0,
			Products: []domain.HouseholdProduct{
				{
					ProductID:  primitive.NewObjectID(),
					CategoryID: c1.CategoryID,
					ImageURL:   f.ImageURL(200, 300),
					Name:       f.StreetName(),
					ISBN:       f.BuzzWord(),
					Settings:   f.LoremIpsumWord(),
					Price:      1223,
					PriceGlob:  1231,
				},
				{
					ProductID:  primitive.NewObjectID(),
					CategoryID: c1.CategoryID,
					ImageURL:   f.ImageURL(300, 300),
					Name:       f.StreetName(),
					ISBN:       f.BuzzWord(),
					Settings:   f.LoremIpsumWord(),
					Price:      12235,
					PriceGlob:  12315,
				},
			},
		})

		err := s.repositories.HouseholdCategory.Save(ctx, c1)
		require.NoError(err)

		p1, p2 := c1.Subcategories[0].Products[0], c1.Subcategories[0].Products[1]
		cart := domain.HouseholdCart([]domain.HouseholdProduct{})
		cart.Add(p1)
		cart.Add(p2)
		customer.Cart = cart

		err = s.repositories.HouseholdCustomer.Save(ctx, customer)
		require.NoError(err)

		// Update c1, remove p2
		c1.Subcategories[0].Products = c1.Subcategories[0].Products[:1]
		err = s.repositories.HouseholdCategory.Update(ctx, c1.CategoryID, dto.UpdateCategoryDTO{
			Subcategories: &c1.Subcategories,
		})
		require.NoError(err)

		// Checks cart
		err = s.tghandler.AskForFIO(ctx, telegramID)
		require.NoError(err)

		require.NoError(s.repositories.HouseholdCustomer.Delete(ctx, customer.CustomerID))
		require.NoError(s.repositories.HouseholdCategory.Delete(ctx, c1.CategoryID))
	})

	s.T().Run("should create order because all products exist at the moment", func(t *testing.T) {
		var (
			username   = f.Username()
			telegramID = f.Int64()
		)
		customer := domain.NewHouseholdCustomer(telegramID, username)

		c1 := domain.NewHouseholdCategory(f.Word(), true)
		c1.CategoryID = primitive.NewObjectID()
		c1.Subcategories = append(c1.Subcategories, domain.Subcategory{
			SubcategoryID: primitive.NewObjectID(),
			Title:         f.Word(),
			Active:        true,
			Rank:          0,
			Products: []domain.HouseholdProduct{
				{
					ProductID:  primitive.NewObjectID(),
					CategoryID: c1.CategoryID,
					ImageURL:   f.ImageURL(200, 300),
					Name:       f.StreetName(),
					ISBN:       f.BuzzWord(),
					Settings:   f.LoremIpsumWord(),
					Price:      1223,
					PriceGlob:  1231,
				},
				{
					ProductID:  primitive.NewObjectID(),
					CategoryID: c1.CategoryID,
					ImageURL:   f.ImageURL(300, 300),
					Name:       f.StreetName(),
					ISBN:       f.BuzzWord(),
					Settings:   f.LoremIpsumWord(),
					Price:      12235,
					PriceGlob:  12315,
				},
			},
		})

		err := s.repositories.HouseholdCategory.Save(ctx, c1)
		require.NoError(err)

		p1, p2 := c1.Subcategories[0].Products[0], c1.Subcategories[0].Products[1]
		cart := domain.HouseholdCart([]domain.HouseholdProduct{})
		cart.Add(p1)
		cart.Add(p2)
		customer.Cart = cart

		err = s.repositories.HouseholdCustomer.Save(ctx, customer)
		require.NoError(err)

		// Checks cart
		s.mockBot = new(MockBot)
		err = s.tghandler.AskForFIO(ctx, telegramID)
		require.NoError(err)

		s.repositories.HouseholdCustomer.Delete(ctx, customer.CustomerID)
		s.repositories.HouseholdCategory.Delete(ctx, c1.CategoryID)
	})
}

func newTgMessage(msgID int, senderTelegramID int64, senderUsername string, text string) *tg.Message {
	return &tg.Message{
		MessageID: msgID,
		From: &tg.User{
			ID:       senderTelegramID,
			IsBot:    false,
			UserName: senderUsername,
		},
		SenderChat: nil,
		Date:       0,
		Chat: &tg.Chat{
			ID: senderTelegramID,
		},
		Text: text,
	}
}
