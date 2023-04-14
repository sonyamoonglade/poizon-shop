package tests

import (
	"context"

	"domain"
	f "github.com/brianvoe/gofakeit/v6"
)

func (s *AppTestSuite) TestHandleOrderTypeInput() {
	var (
		require    = s.Require()
		telegramID = f.Int64()
		username   = f.Username()
		ctx        = context.Background()
	)

	s.Run("valid order type", func() {
		customer := domain.NewHouseholdCustomer(telegramID, username)
		customer.State = domain.StateWaitingForOrderType

		err := s.repositories.HouseholdCustomer.Save(ctx, customer)
		require.NoError(err)

		args := []string{domain.OrderTypeExpress.String()}
		err = s.tghandler.HandleOrderTypeInput(ctx, telegramID, args)
		require.NoError(err)

		dbCustomer, err := s.repositories.HouseholdCustomer.GetByTelegramID(ctx, telegramID)
		require.NoError(err)

		require.NotNil(dbCustomer.Meta.NextOrderType)
		require.Equal(domain.OrderTypeExpress, *dbCustomer.Meta.NextOrderType)

		s.repositories.HouseholdCustomer.Delete(ctx, dbCustomer.CustomerID)
	})

}
