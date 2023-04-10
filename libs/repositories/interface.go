package repositories

import (
	"context"

	"domain"
	"dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Customer interface {
	GetByTelegramID(ctx context.Context, telegramID int64) (domain.Customer, error)
	All(ctx context.Context) ([]domain.Customer, error)
	Save(ctx context.Context, c domain.Customer) error
	UpdateState(ctx context.Context, telegramID int64, newState domain.State) error
	NullifyCatalogOffsets(ctx context.Context) error
	Update(ctx context.Context, customerID primitive.ObjectID, dto dto.UpdateClothingCustomerDTO) error
	Delete(ctx context.Context, customerID primitive.ObjectID) error
}

type Order interface {
	GetByShortID(ctx context.Context, shortID string) (domain.ClothingOrder, error)
	GetAllForCustomer(ctx context.Context, customerID primitive.ObjectID) ([]domain.ClothingOrder, error)
	GetAll(ctx context.Context) ([]domain.ClothingOrder, error)
	Save(ctx context.Context, o domain.ClothingOrder) error
	Approve(ctx context.Context, orderID primitive.ObjectID) (domain.ClothingOrder, error)
	AddComment(ctx context.Context, dto dto.AddCommentDTO) (domain.ClothingOrder, error)
	UpdateToPaid(ctx context.Context, customerID primitive.ObjectID, shortID string) error
	ChangeStatus(ctx context.Context, dto dto.ChangeOrderStatusDTO) (domain.ClothingOrder, error)
	Delete(ctx context.Context, orderID primitive.ObjectID) error
}

type ClothingCatalog interface {
	GetCatalog(ctx context.Context) ([]domain.ClothingProduct, error)
	GetIDByRank(ctx context.Context, rank uint) (primitive.ObjectID, error)
	GetRankByID(ctx context.Context, itemID primitive.ObjectID) (uint, error)
	GetLastRank(ctx context.Context) (uint, error)
	AddItem(ctx context.Context, item domain.ClothingProduct) error
	RemoveItem(ctx context.Context, itemID primitive.ObjectID) error
}

type HouseholdCategory interface {
	GetByID(ctx context.Context, categoryID primitive.ObjectID) (domain.HouseholdCategory, error)
	GetByTitle(ctx context.Context, title string) (domain.HouseholdCategory, error)
	GetProductsByCategoryAndSubcategory(ctx context.Context, cTitle, sTitle string, availableInStock bool) ([]domain.HouseholdProduct, error)
	GetAll(ctx context.Context) ([]domain.HouseholdCategory, error)
	Save(ctx context.Context, c domain.HouseholdCategory) error
	GetTopRank(ctx context.Context) (uint, error)
	Delete(ctx context.Context, categoryID primitive.ObjectID) error
	Update(ctx context.Context, categoryID primitive.ObjectID, dto dto.UpdateCategoryDTO) error
}
