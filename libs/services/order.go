package services

import (
	"context"
	"errors"
	"fmt"

	"domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"nanoid"
	"repositories"
)

type orderService struct {
	repo repositories.Order
}

func NewOrderService(repo repositories.Order) *orderService {
	return &orderService{
		repo: repo,
	}
}
func (o *orderService) GetFreeShortID(ctx context.Context) (string, error) {
	for {
		shortID := nanoid.GenerateNanoID()
		_, err := o.repo.GetByShortID(ctx, shortID)
		if err != nil {
			if errors.Is(err, domain.ErrOrderNotFound) {
				return shortID, nil
			}
			return "", fmt.Errorf("get by short id: %w", err)
		}
		// if reached, means something has been found - skip and go again
		continue
	}
}

func (o *orderService) Save(ctx context.Context, order domain.ClothingOrder) error {
	return o.repo.Save(ctx, order)
}

func (o *orderService) UpdateToPaid(ctx context.Context, customerID primitive.ObjectID, shortID string) error {
	return o.repo.UpdateToPaid(ctx, customerID, shortID)
}

func (o *orderService) GetAllForCustomer(ctx context.Context, customerID primitive.ObjectID) ([]domain.ClothingOrder, error) {
	return o.repo.GetAllForCustomer(ctx, customerID)
}
