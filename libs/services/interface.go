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
	GetLast(ctx context.Context, customerID primitive.ObjectID) (T, error)
	Save(ctx context.Context, order T) error
	UpdateToPaid(ctx context.Context, customerID primitive.ObjectID, shortID string) error
	GetFreeShortID(ctx context.Context) (string, error)
	HasOnlyOneOrder(ctx context.Context, customerID primitive.ObjectID) (bool, error)
}

type HouseholdCategory interface {
	GetByID(ctx context.Context, categoryID primitive.ObjectID) (domain.HouseholdCategory, error)
	GetAll(ctx context.Context) ([]domain.HouseholdCategory, error)
	GetByTitle(ctx context.Context, title string, inStock bool) (domain.HouseholdCategory, error)
	GetAllByInStock(ctx context.Context, inStock bool) ([]domain.HouseholdCategory, error)
	GetProductsByCategoryAndSubcategory(ctx context.Context, cTitle, sTitle string, inStock bool) ([]domain.HouseholdProduct, error)
	New(ctx context.Context, title string, inStock bool) error
	Delete(ctx context.Context, categoryID primitive.ObjectID) error
	Update(ctx context.Context, categoryID primitive.ObjectID, dto dto.UpdateCategoryDTO) error
	CheckIfAllProductsExist(ctx context.Context, cart []domain.HouseholdProduct, inStock bool) (bool, domain.HouseholdProduct, error)
}

type Deleter interface {
	DeleteFromCatalog(c telegram.CatalogMsg) error
}

type HouseholdCatalogMsg interface {
	GetAll(ctx context.Context) ([]telegram.CatalogMsg, error)
	Save(ctx context.Context, m telegram.CatalogMsg) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	DeleteByMsgID(ctx context.Context, msgID int) error
	WipeAll(ctx context.Context, catalogDeleter Deleter) error
}
type customerConstraint interface {
	domain.HouseholdCustomer | domain.ClothingCustomer
}
type dtoConstraint interface {
	dto.UpdateClothingCustomerDTO | dto.UpdateHouseholdCustomerDTO
}
type Customer[T customerConstraint, D dtoConstraint] interface {
	GetAllByPromocodeID(ctx context.Context, promocodeID primitive.ObjectID) ([]T, error)
	GetByTelegramID(ctx context.Context, telegramID int64) (T, error)
	All(ctx context.Context) ([]T, error)
	Save(ctx context.Context, c T) error
	UpdateState(ctx context.Context, telegramID int64, newState domain.State) error
	Update(ctx context.Context, customerID primitive.ObjectID, dto D) error
	Delete(ctx context.Context, customerID primitive.ObjectID) error
}

type ClothingCustomer interface {
	Customer[domain.ClothingCustomer, dto.UpdateClothingCustomerDTO]
	NullifyCatalogOffsets(ctx context.Context) error
}

type HouseholdCustomer interface {
	Customer[domain.HouseholdCustomer, dto.UpdateHouseholdCustomerDTO]
}
