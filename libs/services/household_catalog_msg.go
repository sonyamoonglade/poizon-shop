package services

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"household_bot/pkg/telegram"
	"logger"
	"onlineshop/database"
	"repositories"
)

type householdCatalogMsgService struct {
	repo       repositories.HouseholdCatalogMsg
	transactor database.Transactor
}

func NewHouseholdCatalogMsgService(repo repositories.HouseholdCatalogMsg, transactor database.Transactor) *householdCatalogMsgService {
	return &householdCatalogMsgService{
		repo:       repo,
		transactor: transactor,
	}
}

func (h householdCatalogMsgService) Save(ctx context.Context, m telegram.CatalogMsg) error {
	return h.repo.Save(ctx, m)
}

func (h householdCatalogMsgService) GetAll(ctx context.Context) ([]telegram.CatalogMsg, error) {
	return h.repo.GetAll(ctx)
}

func (h householdCatalogMsgService) Delete(ctx context.Context, id primitive.ObjectID) error {
	return h.repo.Delete(ctx, id)

}

func (h householdCatalogMsgService) DeleteByMsgID(ctx context.Context, msgID int) error {
	return h.repo.DeleteByMsgID(ctx, msgID)
}

// todo: test
func (h householdCatalogMsgService) WipeAll(ctx context.Context, catalogDeleter Deleter) error {
	return h.transactor.WithTransaction(ctx, func(tx context.Context) error {
		msgs, err := h.GetAll(tx)
		if err != nil {
			return fmt.Errorf("get all: %w", err)
		}

		// var err error
		var errors error
		sem := make(chan struct{}, 5)

		for _, m := range msgs {
			sem <- struct{}{}
			// shadowing
			m := m
			go func() {
				defer func() {
					<-sem
				}()
				err := catalogDeleter.DeleteFromCatalog(m)
				// Dont delete from db
				if err != nil {
					const msgNotFound = "message to delete not found"
					// Fine case for us, just delete in DB
					if strings.Contains(err.Error(), msgNotFound) {
						if err := h.repo.Delete(ctx, m.ID); err != nil {
							errors = multierr.Append(errors, fmt.Errorf("delete: %w", err))
						}
						return
					}
					logger.Get().Error("catalog deleter dit not delete msg", zap.Int("msgID", m.MsgID), zap.Error(err))
					return
				}

				if err := h.repo.Delete(ctx, m.ID); err != nil {
					errors = multierr.Append(errors, fmt.Errorf("delete: %w", err))
				}
			}()
		}

		return errors
	})
}
