package tests

import (
	"context"
	"strconv"
	"strings"

	"domain"
	f "github.com/brianvoe/gofakeit/v6"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/mock"
)

func (s *AppTestSuite) TestAskForCalculatorOrderType() {
	var (
		require    = s.Require()
		telegramID = f.Int64()
		username   = f.Username()
		ctx        = context.Background()
	)
	customer := domain.NewCustomer(telegramID, username)
	err := s.repositories.ClothingCustomer.Save(ctx, customer)
	require.NoError(err)

	err = s.tghandler.AskForCalculatorOrderType(ctx, telegramID)
	require.NoError(err)

	dbCustomer, err := s.repositories.ClothingCustomer.GetByTelegramID(ctx, telegramID)
	require.NoError(err)

	require.Equal(domain.StateWaitingForCalculatorOrderType, dbCustomer.TgState)
	s.repositories.ClothingCustomer.Delete(ctx, dbCustomer.CustomerID)
}
func (s *AppTestSuite) TestHandleCalculatorOrderTypeInput() {
	var (
		require    = s.Require()
		telegramID = f.Int64()
		username   = f.Username()
		ctx        = context.Background()
	)
	for _, ordTyp := range []domain.OrderType{domain.OrderTypeExpress, domain.OrderTypeNormal} {
		customer := domain.NewCustomer(telegramID, username)
		// Set appropriate state and category
		customer.TgState = domain.StateWaitingForCalculatorOrderType
		customer.UpdateLastEditPositionCategory(domain.CategoryOther)

		err := s.repositories.ClothingCustomer.Save(ctx, customer)
		require.NoError(err)

		err = s.tghandler.HandleCalculatorOrderTypeInput(ctx, telegramID, ordTyp)
		require.NoError(err)

		dbCustomer, err := s.repositories.ClothingCustomer.GetByTelegramID(ctx, telegramID)
		require.NoError(err)

		require.NotNil(dbCustomer.CalculatorMeta)
		require.Equal(domain.StateWaitingForCalculatorCategory, dbCustomer.TgState)
		require.Equal(ordTyp, *dbCustomer.CalculatorMeta.NextOrderType)

		s.repositories.ClothingCustomer.Delete(ctx, dbCustomer.CustomerID)
	}
}

func (s *AppTestSuite) TestAskForCalculatorCategory() {
	var (
		require    = s.Require()
		telegramID = f.Int64()
		username   = f.Username()
		ctx        = context.Background()
	)
	customer := domain.NewCustomer(telegramID, username)
	err := s.repositories.ClothingCustomer.Save(ctx, customer)
	require.NoError(err)

	err = s.tghandler.AskForCalculatorCategory(ctx, telegramID)
	require.NoError(err)

	dbCustomer, err := s.repositories.ClothingCustomer.GetByTelegramID(ctx, telegramID)
	require.NoError(err)

	require.Equal(domain.StateWaitingForCalculatorCategory, dbCustomer.TgState)
	s.repositories.ClothingCustomer.Delete(ctx, dbCustomer.CustomerID)
}

func (s *AppTestSuite) TestHandleCalculatorCategoryInput() {
	var (
		require    = s.Require()
		telegramID = f.Int64()
		username   = f.Username()
		ctx        = context.Background()
	)

	for _, category := range []domain.Category{domain.CategoryHeavy, domain.CategoryOther, domain.CategoryLight} {
		customer := domain.NewCustomer(telegramID, username)
		// Set appropriate state
		customer.TgState = domain.StateWaitingForCalculatorCategory
		customer.UpdateCalculatorMetaOrderType(domain.OrderTypeNormal)

		err := s.repositories.ClothingCustomer.Save(ctx, customer)
		require.NoError(err)

		err = s.tghandler.HandleCalculatorCategoryInput(ctx, telegramID, category)
		require.NoError(err)

		dbCustomer, err := s.repositories.ClothingCustomer.GetByTelegramID(ctx, telegramID)
		require.NoError(err)

		require.Equal(category, *dbCustomer.CalculatorMeta.Category)
		require.Equal(domain.StateWaitingForCalculatorInput, dbCustomer.TgState)

		s.repositories.ClothingCustomer.Delete(ctx, dbCustomer.CustomerID)
	}
}

func (s *AppTestSuite) TestHandleCalculatorInput() {
	var (
		require    = s.Require()
		telegramID = f.Int64()
		username   = f.Username()
		ctx        = context.Background()
	)

	type test struct {
		expected uint64
		meta     domain.CalculatorMeta
		input    string
	}
	var tests []test
	inputs := make([]uint16, 15)
	f.Slice(&inputs)
	for _, typ := range []domain.OrderType{domain.OrderTypeExpress, domain.OrderTypeNormal} {
		for _, cat := range []domain.Category{domain.CategoryHeavy, domain.CategoryLight, domain.CategoryOther} {
			for _, yuanInput := range inputs {
				args := domain.ConvertYuanArgs{
					X:         uint64(yuanInput),
					Rate:      11.8,
					OrderType: typ,
					Category:  cat,
				}
				expected := domain.ConvertYuan(args)
				meta := domain.CalculatorMeta{
					NextOrderType: &typ,
					Category:      &cat,
				}
				tests = append(tests, test{
					expected: expected,
					meta:     meta,
					input:    strconv.Itoa(int(yuanInput)),
				})
			}
		}
	}
	// Custom logic
	s.mockBot = new(MockBot)
	for _, test := range tests {
		s.mockBot.On("Send", mock.Anything).Run(func(args mock.Arguments) {
			msgToTg, ok := args.Get(0).(*tg.Message)
			require.True(ok)
			// Check that outgoing message contains valid price
			ok = strings.Contains(msgToTg.Text, strconv.Itoa(int(test.expected)))
			require.True(ok)
		}).Return(tg.Message{}, nil).Times(1)

		customer := domain.NewCustomer(telegramID, username)
		// Set appropriate state, meta
		customer.TgState = domain.StateWaitingForCalculatorInput
		customer.CalculatorMeta = test.meta
		err := s.repositories.ClothingCustomer.Save(ctx, customer)
		require.NoError(err)

		m := newTgMessage(f.IntRange(1, 10), telegramID, username, test.input)
		err = s.tghandler.HandleCalculatorInput(ctx, m)
		require.NoError(err)

		dbCustomer, err := s.repositories.ClothingCustomer.GetByTelegramID(ctx, telegramID)
		require.NoError(err)

		s.repositories.ClothingCustomer.Delete(ctx, dbCustomer.CustomerID)
	}
}
