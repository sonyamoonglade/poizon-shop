package tests

import (
	"context"
	"strconv"

	"domain"
	f "github.com/brianvoe/gofakeit/v6"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"utils/url"
)

func (s *AppTestSuite) TestAddPosition() {
	var (
		require    = s.Require()
		telegramID = f.Int64()
		username   = f.Username()
		ctx        = context.Background()
	)

	s.Run("customer's cart is empty", func() {
		customer := domain.NewCustomer(telegramID, username)
		err := s.repositories.ClothingCustomer.Save(ctx, customer)
		require.NoError(err)

		// In case customer's cart is empty, so h.askForOrderType will be called
		err = s.tghandler.AddPosition(ctx, telegramID)
		require.NoError(err)

		dbCustomer, err := s.repositories.ClothingCustomer.GetByTelegramID(ctx, telegramID)
		require.NoError(err)

		require.Equal(domain.StateWaitingForOrderType, dbCustomer.TgState)
		// cleanup
		s.repositories.ClothingCustomer.Delete(ctx, dbCustomer.CustomerID)
	})

	s.Run("customer's cart is not empty", func() {
		customer := domain.NewCustomer(telegramID, username)
		customer.Cart.Add(domain.ClothingPosition{
			PositionID: primitive.NewObjectID(),
			ShopLink:   f.URL(),
			PriceRUB:   123,
			PriceYUAN:  245,
			Button:     "a",
			Size:       "b",
			Category:   "c",
		})
		err := s.repositories.ClothingCustomer.Save(ctx, customer)
		require.NoError(err)

		// In case customer's cart is empty, so h.askForCategory will be called
		err = s.tghandler.AddPosition(ctx, telegramID)
		require.NoError(err)

		dbCustomer, err := s.repositories.ClothingCustomer.GetByTelegramID(ctx, telegramID)
		require.NoError(err)

		require.Equal(domain.StateWaitingForCategory, dbCustomer.TgState)

		// cleanup
		s.repositories.ClothingCustomer.Delete(ctx, dbCustomer.CustomerID)
	})

}

func (s *AppTestSuite) TestHandleCategoryInput() {
	var (
		require    = s.Require()
		telegramID = f.Int64()
		username   = f.Username()
		ctx        = context.Background()
	)

	for _, category := range []domain.Category{domain.CategoryHeavy, domain.CategoryOther, domain.CategoryLight} {
		customer := domain.NewCustomer(telegramID, username)
		// Set appropriate state
		customer.TgState = domain.StateWaitingForCategory

		err := s.repositories.ClothingCustomer.Save(ctx, customer)
		require.NoError(err)

		err = s.tghandler.HandleCategoryInput(ctx, telegramID, category)
		require.NoError(err)

		dbCustomer, err := s.repositories.ClothingCustomer.GetByTelegramID(ctx, telegramID)
		require.NoError(err)

		require.NotNil(dbCustomer.LastEditPosition)
		require.Equal(category, dbCustomer.LastEditPosition.Category)
		require.Equal(domain.StateWaitingForSize, dbCustomer.TgState)

		s.repositories.ClothingCustomer.Delete(ctx, dbCustomer.CustomerID)
	}
}

func (s *AppTestSuite) TestHandleOrderTypeInput() {
	var (
		require    = s.Require()
		telegramID = f.Int64()
		username   = f.Username()
		ctx        = context.Background()
	)

	for _, ordTyp := range []domain.OrderType{domain.OrderTypeExpress, domain.OrderTypeNormal} {
		customer := domain.NewCustomer(telegramID, username)
		// Set appropriate state and category
		customer.TgState = domain.StateWaitingForOrderType
		customer.UpdateLastEditPositionCategory(domain.CategoryOther)

		err := s.repositories.ClothingCustomer.Save(ctx, customer)
		require.NoError(err)

		err = s.tghandler.HandleOrderTypeInput(ctx, telegramID, ordTyp)
		require.NoError(err)

		dbCustomer, err := s.repositories.ClothingCustomer.GetByTelegramID(ctx, telegramID)
		require.NoError(err)

		require.NotNil(dbCustomer.LastEditPosition)
		require.Equal(domain.StateWaitingForCategory, dbCustomer.TgState)
		require.Equal(ordTyp, *dbCustomer.Meta.NextOrderType)

		s.repositories.ClothingCustomer.Delete(ctx, dbCustomer.CustomerID)
	}
}

func (s *AppTestSuite) TestHandleSizeInput() {
	var (
		require    = s.Require()
		telegramID = f.Int64()
		username   = f.Username()
		ctx        = context.Background()
	)
	s.mockBot.On("Send", mock.Anything).Return(tg.Message{}, nil)
	sizes := []string{"A", "B", "C", "L", "#"}
	for _, size := range sizes {
		customer := domain.NewCustomer(telegramID, username)
		// Set appropriate state, category, order type
		customer.TgState = domain.StateWaitingForSize
		customer.UpdateLastEditPositionCategory(domain.CategoryOther)
		var t = domain.OrderTypeExpress
		customer.Meta.NextOrderType = &t

		err := s.repositories.ClothingCustomer.Save(ctx, customer)
		require.NoError(err)

		m := newTgMessage(f.IntRange(1, 10), telegramID, username, size)
		err = s.tghandler.HandleSizeInput(ctx, m)
		require.NoError(err)

		dbCustomer, err := s.repositories.ClothingCustomer.GetByTelegramID(ctx, telegramID)
		require.NoError(err)

		require.NotNil(dbCustomer.LastEditPosition)
		require.Equal(dbCustomer.LastEditPosition.Size, size)
		require.Equal(domain.StateWaitingForButton, dbCustomer.TgState)

		s.repositories.ClothingCustomer.Delete(ctx, dbCustomer.CustomerID)
	}
}

func (s *AppTestSuite) HandleButtonSelect() {
	var (
		require    = s.Require()
		telegramID = f.Int64()
		username   = f.Username()
		ctx        = context.Background()
	)

	for _, b := range []domain.Button{domain.Button95, domain.ButtonTorqoise, domain.ButtonGrey} {
		customer := domain.NewCustomer(telegramID, username)
		// Set appropriate state, category, order type, size
		customer.TgState = domain.StateWaitingForButton
		customer.UpdateLastEditPositionCategory(domain.CategoryOther)
		var t = domain.OrderTypeExpress
		customer.Meta.NextOrderType = &t
		customer.LastEditPosition.Size = "L"

		err := s.repositories.ClothingCustomer.Save(ctx, customer)
		require.NoError(err)

		c := newCallback(f.IntRange(1, 10), telegramID, username, "0")
		err = s.tghandler.HandleButtonSelect(ctx, c, b)
		require.NoError(err)

		dbCustomer, err := s.repositories.ClothingCustomer.GetByTelegramID(ctx, telegramID)
		require.NoError(err)

		require.NotNil(dbCustomer.LastEditPosition)
		require.Equal(dbCustomer.LastEditPosition.Button, b)
		require.Equal(domain.StateWaitingForPrice, dbCustomer.TgState)

		s.repositories.ClothingCustomer.Delete(ctx, dbCustomer.CustomerID)
	}
}

func (s *AppTestSuite) HandleTestPriceInput() {
	var (
		require    = s.Require()
		telegramID = f.Int64()
		username   = f.Username()
		ctx        = context.Background()
	)

	inputs := make([]int, 10)
	f.Slice(&inputs)
	for _, inputYuan := range inputs {
		customer := domain.NewCustomer(telegramID, username)
		// Set appropriate state, category, order type, size, button
		customer.TgState = domain.StateWaitingForPrice
		customer.LastEditPosition.Category = domain.CategoryOther
		var t = domain.OrderTypeExpress
		customer.Meta.NextOrderType = &t
		customer.LastEditPosition.Size = "L"
		customer.LastEditPosition.Button = domain.ButtonTorqoise

		err := s.repositories.ClothingCustomer.Save(ctx, customer)
		require.NoError(err)

		inpStr := strconv.Itoa(inputYuan)
		m := newTgMessage(f.IntRange(1, 10), telegramID, username, inpStr)
		err = s.tghandler.HandlePriceInput(ctx, m)
		require.NoError(err)

		dbCustomer, err := s.repositories.ClothingCustomer.GetByTelegramID(ctx, telegramID)
		require.NoError(err)

		require.NotNil(dbCustomer.LastEditPosition)
		expectedPriceRub := domain.ConvertYuan(domain.ConvertYuanArgs{
			X:         uint64(inputYuan),
			Rate:      11.8,
			OrderType: *customer.Meta.NextOrderType,
			Category:  customer.LastEditPosition.Category,
		})
		require.Equal(dbCustomer.LastEditPosition.PriceRUB, expectedPriceRub)
		require.Equal(domain.StateWaitingForLink, dbCustomer.TgState)

		s.repositories.ClothingCustomer.Delete(ctx, dbCustomer.CustomerID)
	}
}

func (s *AppTestSuite) TestHandleLinkInput() {
	var (
		require    = s.Require()
		telegramID = f.Int64()
		username   = f.Username()
		ctx        = context.Background()
	)

	links := []string{"https://dw4.co/t/A/abdc", "https://dw", "https://google.com"}

	for _, link := range links {
		customer := domain.NewCustomer(telegramID, username)
		// Set appropriate state, category, order type, size, button
		customer.TgState = domain.StateWaitingForLink
		customer.UpdateLastEditPositionCategory(domain.CategoryOther)
		var t = domain.OrderTypeExpress
		customer.Meta.NextOrderType = &t
		customer.LastEditPosition.Size = "L"
		customer.LastEditPosition.Button = domain.ButtonTorqoise
		customer.LastEditPosition.PriceRUB = 100
		customer.LastEditPosition.PriceYUAN = 50

		err := s.repositories.ClothingCustomer.Save(ctx, customer)
		require.NoError(err)

		m := newTgMessage(f.IntRange(1, 10), telegramID, username, link)
		err = s.tghandler.HandleLinkInput(ctx, m)
		require.NoError(err)

		dbCustomer, err := s.repositories.ClothingCustomer.GetByTelegramID(ctx, telegramID)
		require.NoError(err)

		s.repositories.ClothingCustomer.Delete(ctx, dbCustomer.CustomerID)

		if !url.IsValidDW4URL(link) {
			require.NotNil(dbCustomer.LastEditPosition)
			require.Zero(dbCustomer.LastEditPosition.ShopLink)
			require.Equal(domain.StateWaitingForLink, dbCustomer.TgState)
			continue
		}

		require.Equal(dbCustomer.LastEditPosition.ShopLink, link)
		require.Equal(domain.StateDefault, dbCustomer.TgState)
	}
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

func newCallback(msgID int, senderTelegramID int64, senderUsername string, data string) *tg.CallbackQuery {
	return &tg.CallbackQuery{
		From: &tg.User{
			ID:       senderTelegramID,
			UserName: senderUsername,
		},
		Message: &tg.Message{
			MessageID: msgID,
		},
		Data: data,
	}
}
