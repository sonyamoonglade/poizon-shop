package services

import (
	"context"

	"domain"
	"dto"
	"household_bot/pkg/telegram"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order[T domain.HouseholdOrder | domain.ClothingOrder] interface {
	GetAllForCustomer(ctx context.Context, customerID primitive.ObjectID) ([]T, error)
	Save(ctx context.Context, order T) error
	UpdateToPaid(ctx context.Context, customerID primitive.ObjectID, shortID string) error
	GetFreeShortID(ctx context.Context) (string, error)
}

type HouseholdCategory interface {
	GetByID(ctx context.Context, categoryID primitive.ObjectID) (domain.HouseholdCategory, error)
	GetAll(ctx context.Context) ([]domain.HouseholdCategory, error)
	GetAllByInStock(ctx context.Context, inStock bool) ([]domain.HouseholdCategory, error)
	New(ctx context.Context, title string, inStock bool) error
	Delete(ctx context.Context, categoryID primitive.ObjectID) error
	Update(ctx context.Context, categoryID primitive.ObjectID, dto dto.UpdateCategoryDTO) error
	CheckIfAllProductsExist(ctx context.Context, productIDs []primitive.ObjectID) (bool, error)
}

type Deleter interface {
	DeleteFromCatalog(c telegram.CatalogMsg) error
}

type HouseholdCatalogMsg interface {
	Save(ctx context.Context, m telegram.CatalogMsg) error
	GetAll(ctx context.Context) ([]telegram.CatalogMsg, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
	DeleteByMsgID(ctx context.Context, msgID int) error
	WipeAll(ctx context.Context, catalogDeleter Deleter) error
}
