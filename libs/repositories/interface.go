package repositories

import (
	"context"

	"domain"
	"dto"
	"household_bot/pkg/telegram"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// todo: unify to generic customer interface!!
type ClothingCustomer interface {
	GetByTelegramID(ctx context.Context, telegramID int64) (domain.ClothingCustomer, error)
	GetAllByPromocodeID(ctx context.Context, promocodeID primitive.ObjectID) ([]domain.ClothingCustomer, error)
	All(ctx context.Context) ([]domain.ClothingCustomer, error)
	Save(ctx context.Context, c domain.ClothingCustomer) error
	UpdateState(ctx context.Context, telegramID int64, newState domain.State) error
	Update(ctx context.Context, customerID primitive.ObjectID, dto dto.UpdateClothingCustomerDTO) error
	Delete(ctx context.Context, customerID primitive.ObjectID) error

	NullifyCatalogOffsets(ctx context.Context) error
}

// todo: unify to generic customer interface!!
type HouseholdCustomer interface {
	GetByTelegramID(ctx context.Context, telegramID int64) (domain.HouseholdCustomer, error)
	GetAllByPromocodeID(ctx context.Context, promocodeID primitive.ObjectID) ([]domain.ClothingCustomer, error)
	All(ctx context.Context) ([]domain.HouseholdCustomer, error)
	Save(ctx context.Context, c domain.HouseholdCustomer) error
	UpdateState(ctx context.Context, telegramID int64, newState domain.State) error
	Update(ctx context.Context, customerID primitive.ObjectID, dto dto.UpdateHouseholdCustomerDTO) error
	Delete(ctx context.Context, customerID primitive.ObjectID) error
}

type Order[T domain.ClothingOrder | domain.HouseholdOrder] interface {
	GetByShortID(ctx context.Context, shortID string) (T, error)
	GetAllForCustomer(ctx context.Context, customerID primitive.ObjectID) ([]T, error)
	GetLast(ctx context.Context, customerID primitive.ObjectID) ([]T, error)
	GetAll(ctx context.Context) ([]T, error)
	Save(ctx context.Context, o T) error
	Approve(ctx context.Context, orderID primitive.ObjectID) (T, error)
	AddComment(ctx context.Context, dto dto.AddCommentDTO) (T, error)
	UpdateToPaid(ctx context.Context, customerID primitive.ObjectID, shortID string) error
	ChangeStatus(ctx context.Context, dto dto.ChangeOrderStatusDTO) (T, error)
	Delete(ctx context.Context, orderID primitive.ObjectID) error
}

type ClothingCatalog interface {
	GetCatalog(ctx context.Context) ([]domain.ClothingProduct, error)
	GetIDByRank(ctx context.Context, rank uint) (primitive.ObjectID, error)
	GetRankByID(ctx context.Context, itemID primitive.ObjectID) (uint, error)
	GetTopRank(ctx context.Context) (uint, error)
	AddItem(ctx context.Context, item domain.ClothingProduct) error
	RemoveItem(ctx context.Context, itemID primitive.ObjectID) error
}
type HouseholdCategory interface {
	GetByID(ctx context.Context, categoryID primitive.ObjectID) (domain.HouseholdCategory, error)
	GetByTitle(ctx context.Context, title string, inStock bool) (domain.HouseholdCategory, error)
	GetProductsByCategoryAndSubcategory(ctx context.Context, cTitle, sTitle string, inStock bool) ([]domain.HouseholdProduct, error)
	GetAll(ctx context.Context) ([]domain.HouseholdCategory, error)
	GetAllByInStock(ctx context.Context, inStock bool) ([]domain.HouseholdCategory, error)
	Save(ctx context.Context, c domain.HouseholdCategory) error
	GetTopRank(ctx context.Context, inStock bool) (uint, error)
	Delete(ctx context.Context, categoryID primitive.ObjectID) error
	Update(ctx context.Context, categoryID primitive.ObjectID, dto dto.UpdateCategoryDTO) error
}

type HouseholdCatalogMsg interface {
	GetAll(ctx context.Context) ([]telegram.CatalogMsg, error)
	Save(ctx context.Context, m telegram.CatalogMsg) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	DeleteByMsgID(ctx context.Context, msgID int) error
}

type Promocode interface {
	GetAll(ctx context.Context) ([]domain.Promocode, error)
	GetByID(ctx context.Context, promocodeID primitive.ObjectID) (domain.Promocode, error)
	GetByShortID(ctx context.Context, shortID string) (domain.Promocode, error)
	Save(ctx context.Context, promo domain.Promocode) error
	Delete(ctx context.Context, promocodeID primitive.ObjectID) error
	Update(ctx context.Context, promocodeID primitive.ObjectID, dto dto.UpdatePromocodeDTO) error
}
