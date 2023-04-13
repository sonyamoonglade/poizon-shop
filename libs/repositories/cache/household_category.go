package cache

import (
	"context"
	"domain"
	"dto"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *redisCache) GetByID(ctx context.Context, categoryID primitive.ObjectID) (domain.HouseholdCategory, error) {
	panic("not implemented") // TODO: Implement
}

func (r *redisCache) GetByTitle(ctx context.Context, title string) (domain.HouseholdCategory, error) {
	panic("not implemented") // TODO: Implement
}

func (r *redisCache) GetProductsByCategoryAndSubcategory(ctx context.Context, cTitle string, sTitle string, availableInStock bool) ([]domain.HouseholdProduct, error) {
	panic("not implemented") // TODO: Implement
}

func (r *redisCache) GetAll(ctx context.Context) ([]domain.HouseholdCategory, error) {
	panic("not implemented") // TODO: Implement
}

func (r *redisCache) Save(ctx context.Context, c domain.HouseholdCategory) error {
	panic("not implemented") // TODO: Implement
}

func (r *redisCache) Delete(ctx context.Context, categoryID primitive.ObjectID) error {
	panic("not implemented") // TODO: Implement
}

func (r *redisCache) Update(ctx context.Context, categoryID primitive.ObjectID, dto dto.UpdateCategoryDTO) error {
	panic("not implemented") // TODO: Implement
}

func (r *redisCache) GetTopRank(ctx context.Context) (uint, error) {
	panic("not implemented") // TODO: Implement
}
