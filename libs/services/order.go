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

type orderService[T domain.HouseholdOrder | domain.ClothingOrder] struct {
	repo repositories.Order[T]
}

func NewClothingOrderService(repo repositories.Order[domain.ClothingOrder]) *orderService[domain.ClothingOrder] {
	return &orderService[domain.ClothingOrder]{
		repo: repo,
	}
}

func NewHouseholdOrderService(repo repositories.Order[domain.HouseholdOrder]) *orderService[domain.HouseholdOrder] {
	return &orderService[domain.HouseholdOrder]{
		repo: repo,
	}
}

func (o *orderService[T]) GetFreeShortID(ctx context.Context) (string, error) {
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

func (o *orderService[T]) Save(ctx context.Context, order T) error {
	return o.repo.Save(ctx, order)
}

func (o *orderService[T]) UpdateToPaid(ctx context.Context, customerID primitive.ObjectID, shortID string) error {
	return o.repo.UpdateToPaid(ctx, customerID, shortID)
}

func (o *orderService[T]) GetAllForCustomer(ctx context.Context, customerID primitive.ObjectID) ([]T, error) {
	return o.repo.GetAllForCustomer(ctx, customerID)
}
