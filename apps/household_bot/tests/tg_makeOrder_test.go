package tests

import (
	"context"
	"testing"

	"domain"
	"dto"
	f "github.com/brianvoe/gofakeit/v6"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"household_bot/internal/telegram/templates"
	mock_handler "household_bot/mocks"
)

func (s *AppTestSuite) TestHandlerMakeOrder() {
	var (
		require = s.Require()
		ctx     = context.Background()
		t       = s.T()
	)

	t.Run("should check cart and return product not found message", func(t *testing.T) {
		var (
			username   = f.Username()
			telegramID = f.Int64()
		)

		customer := domain.NewHouseholdCustomer(telegramID, username)
		customer.CustomerID = primitive.NewObjectID()

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
		customer.Cart.Add(p1)
		customer.Cart.Add(p2)

		err = s.repositories.HouseholdCustomer.Save(ctx, customer)
		require.NoError(err)

		// Update c1, remove p2 (p2 is missing in cart)
		c1.Subcategories[0].Products = c1.Subcategories[0].Products[:1]
		err = s.repositories.HouseholdCategory.Update(ctx, c1.CategoryID, dto.UpdateCategoryDTO{
			Subcategories: &c1.Subcategories,
		})
		require.NoError(err)

		// Mocking bot
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mb := mock_handler.NewMockBot(ctrl)
		// For initial message

		mb.EXPECT().
			Send(MsgConfigMatcher(templates.CheckingCart())).
			Return(tg.Message{}, nil).
			Times(1)

		// For product not found
		mb.EXPECT().
			Send(MsgConfigMatcher(templates.ProductNotFound(p2.Name, p2.ISBN))).
			Do(func(args any) {
				m := args.(tg.MessageConfig)
				require.Equal(templates.ProductNotFound(p2.Name, p2.ISBN), m.Text)
			}).
			Return(tg.Message{}, nil).
			Times(1)
		s.replaceBotInHandler(mb)
		// --

		t.Cleanup(func() {
			require.NoError(s.repositories.HouseholdCustomer.Delete(ctx, customer.CustomerID))
			require.NoError(s.repositories.HouseholdCategory.Delete(ctx, c1.CategoryID))
		})

		err = s.tghandler.AskForFIO(ctx, telegramID)
		require.NoError(err)

	})

	t.Run("should check cart and return correct message because cart is fine", func(t *testing.T) {

		var (
			username   = f.Username()
			telegramID = f.Int64()
		)

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

		customer := domain.NewHouseholdCustomer(telegramID, username)
		p1, p2 := c1.Subcategories[0].Products[0], c1.Subcategories[0].Products[1]
		customer.Cart.Add(p1)
		customer.Cart.Add(p2)

		err = s.repositories.HouseholdCustomer.Save(ctx, customer)
		require.NoError(err)

		// Mocking bot
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mb := mock_handler.NewMockBot(ctrl)
		mb.EXPECT().
			Send(gomock.Any()).
			Return(tg.Message{}, nil).
			Times(3)
		s.replaceBotInHandler(mb)
		// --

		t.Cleanup(func() {
			s.repositories.HouseholdCustomer.Delete(ctx, customer.CustomerID)
			s.repositories.HouseholdCategory.Delete(ctx, c1.CategoryID)
		})

		err = s.tghandler.AskForFIO(ctx, telegramID)
		require.NoError(err)
	})

	t.Run("should make order without discount because customer has no promocode", func(t *testing.T) {

		var (
			username   = f.Username()
			telegramID = f.Int64()
		)

		c1 := domain.NewHouseholdCategory(f.Word(), true)
		c1.CategoryID = primitive.NewObjectID()
		c1.Subcategories = append(c1.Subcategories, domain.Subcategory{
			SubcategoryID: primitive.NewObjectID(),
			Title:         f.Word(),
			Active:        true,
			Rank:          0,
			Products: []domain.HouseholdProduct{
				{
					ProductID:   primitive.NewObjectID(),
					CategoryID:  c1.CategoryID,
					Name:        f.StreetName(),
					ISBN:        f.BuzzWord(),
					AvailableIn: &[]string{f.StreetName()},
					Price:       1223,
					PriceGlob:   1231,
				},
				{
					ProductID:   primitive.NewObjectID(),
					CategoryID:  c1.CategoryID,
					Name:        f.StreetName(),
					ISBN:        f.BuzzWord(),
					AvailableIn: &[]string{f.StreetName()},
					Price:       12235,
					PriceGlob:   12315,
				},
			},
		})

		require.NoError(s.repositories.HouseholdCategory.Save(ctx, c1))

		customer := domain.NewHouseholdCustomer(telegramID, username)
		customer.
			UpdateState(domain.StateWaitingForDeliveryAddress).
			SetFullName(f.Name()).
			SetPhoneNumber("89128519000")

		p1, p2 := c1.Subcategories[0].Products[0], c1.Subcategories[0].Products[1]
		customer.Cart.Add(p1)
		customer.Cart.Add(p2)
		cartTotal := p1.Price + p2.Price

		require.NoError(s.repositories.HouseholdCustomer.Save(ctx, customer))
		// Mocking bot
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mb := mock_handler.NewMockBot(ctrl)
		mb.EXPECT().
			Send(gomock.Any()).
			Return(tg.Message{}, nil).
			Times(3)
		s.replaceBotInHandler(mb)
		// --

		t.Cleanup(func() {
			s.repositories.HouseholdCustomer.Delete(ctx, customer.CustomerID)
			s.repositories.HouseholdCategory.Delete(ctx, c1.CategoryID)
		})

		err := s.tghandler.HandleDeliveryAddressInput(
			ctx,
			newTgMessage(
				f.IntRange(1, 200),
				telegramID,
				username,
				f.Address().Address,
			),
		)
		require.NoError(err)
		created, err := s.svc.HouseholdOrder.GetLast(ctx, customer.CustomerID)
		require.NotZero(created)

		require.Equal(cartTotal, created.AmountRUB)
		require.Equal(cartTotal, created.DiscountedAmount)
		require.False(created.IsPaid)
		require.False(created.IsApproved)
		require.Equal(customer.Cart, created.Cart)

	})

	t.Run("should make order with discount because customer has promocode", func(t *testing.T) {
		var (
			username   = f.Username()
			telegramID = f.Int64()
		)

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
					Name:       f.StreetName(),
					ISBN:       f.BuzzWord(),
					Price:      1223,

					AvailableIn: &[]string{f.StreetName()},
					PriceGlob:   1231,
				},
				{
					ProductID:  primitive.NewObjectID(),
					CategoryID: c1.CategoryID,
					Name:       f.StreetName(),
					ISBN:       f.BuzzWord(),
					Price:      12235,

					AvailableIn: &[]string{f.StreetName()},
					PriceGlob:   12315,
				},
			},
		})
		require.NoError(s.repositories.HouseholdCategory.Save(ctx, c1))

		promo := domain.NewPromocode(
			"testing promo",
			domain.DiscountMap{
				domain.SourceHousehold: 200,
				domain.SourceClothing:  100,
			},
			f.Adjective(),
		)
		require.NoError(s.repositories.Promocode.Save(ctx, promo))

		customer := domain.NewHouseholdCustomer(telegramID, username)
		customer.
			UpdateState(domain.StateWaitingForDeliveryAddress).
			SetFullName(f.Name()).
			SetPhoneNumber("89128519000").
			UsePromocode(promo)

		discount := customer.
			MustGetPromocode().
			GetHouseholdDiscount()

		p1, p2 := c1.Subcategories[0].Products[0], c1.Subcategories[0].Products[1]
		customer.Cart.Add(p1)
		customer.Cart.Add(p2)
		cartTotal := (p1.Price - discount) + (p2.Price - discount)
		require.NoError(s.repositories.HouseholdCustomer.Save(ctx, customer))

		// Mocking bot
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mb := mock_handler.NewMockBot(ctrl)
		mb.EXPECT().
			Send(gomock.Any()).
			Return(tg.Message{}, nil).
			Times(3)
		s.replaceBotInHandler(mb)
		// --

		t.Cleanup(func() {
			s.repositories.HouseholdCustomer.Delete(ctx, customer.CustomerID)
			s.repositories.HouseholdCategory.Delete(ctx, c1.CategoryID)
		})

		err := s.tghandler.HandleDeliveryAddressInput(
			ctx,
			newTgMessage(
				f.IntRange(1, 200),
				telegramID,
				username,
				f.Address().Address,
			),
		)
		require.NoError(err)

		created, err := s.svc.HouseholdOrder.GetLast(ctx, customer.CustomerID)
		require.NoError(err)
		require.NotZero(created)

		require.NotEqual(cartTotal, created.AmountRUB)
		require.Equal(cartTotal, created.DiscountedAmount)
		require.True(created.AmountRUB-discount*2 == created.DiscountedAmount)
		require.False(created.IsPaid)
		require.False(created.IsApproved)
		require.Equal(customer.Cart, created.Cart)

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
