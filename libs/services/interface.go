package services

import (
	"context"

	"domain"
	"dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order interface {
	GetAllForCustomer(ctx context.Context, customerID primitive.ObjectID) ([]domain.ClothingOrder, error)
	Save(ctx context.Context, order domain.ClothingOrder) error
	UpdateToPaid(ctx context.Context, customerID primitive.ObjectID, shortID string) error
	GetFreeShortID(ctx context.Context) (string, error)
}

type HouseholdCategory interface {
	New(ctx context.Context, title string) error
	Delete(ctx context.Context, categoryID primitive.ObjectID) error
	Update(ctx context.Context, categoryID primitive.ObjectID, dto dto.UpdateCategoryDTO) error
	GetAll(ctx context.Context) ([]domain.Category, error)
}
