package repositories

import (
	"context"
	"fmt"

	"domain"
	"dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HouseholdOnChangeFunc func(items []domain.HouseholdCategory)

type householdCategoryRepo struct {
	catalog      *mongo.Collection
	onChangeHook HouseholdOnChangeFunc
}

func NewHouseholdCategoryRepo(catalog *mongo.Collection, onChangeHook HouseholdOnChangeFunc) *householdCategoryRepo {
	return &householdCategoryRepo{
		catalog:      catalog,
		onChangeHook: onChangeHook,
	}
}

func (r *householdCategoryRepo) GetTopRank(ctx context.Context) (uint, error) {
	numDocuments, err := r.catalog.CountDocuments(ctx, bson.D{})
	if err != nil {
		return 0, fmt.Errorf("count documents: %w", err)
	}
	return uint(numDocuments), nil
}

func (r *householdCategoryRepo) New(ctx context.Context, c domain.HouseholdCategory) error {
	_, err := r.catalog.InsertOne(ctx, c)
	if err != nil {
		return fmt.Errorf("insert one: %w", err)
	}
	return nil
}

func (r *householdCategoryRepo) Update(ctx context.Context, categoryID primitive.ObjectID, dto dto.UpdateCategoryDTO) error {
	upd := bson.M{}
	fil := bson.M{"_id": categoryID}
	opts := options.Update()
	opts.SetUpsert(true)
	if dto.Active != nil {
		upd["active"] = dto.Active
	}

	if dto.Title != nil {
		upd["title"] = dto.Title
	}

	if dto.Subcategories != nil {
		for is, subcategory := range *dto.Subcategories {
			if subcategory.SubcategoryID.IsZero() {
				(*dto.Subcategories)[is].SubcategoryID = primitive.NewObjectID()
			}
			for ip, product := range subcategory.Products {
				if product.ProductID.IsZero() {
					subcategory.Products[ip].ProductID = primitive.NewObjectID()
				}
			}
		}
		upd["subcategories"] = dto.Subcategories
	}

	if dto.Rank != nil {
		upd["rank"] = dto.Rank
	}

	res, err := r.catalog.UpdateOne(ctx, fil, bson.M{"$set": upd})
	if err != nil {
		return fmt.Errorf("update one: %w", err)
	}
	if res.ModifiedCount == 0 {
		return domain.ErrCategoryNotFound
	}

	return nil
}

func (r *householdCategoryRepo) Delete(ctx context.Context, categoryID primitive.ObjectID) error {
	//TODO implement me
	panic("implement me")
}

func (r *householdCategoryRepo) GetAll(ctx context.Context) ([]domain.Category, error) {
	//TODO implement me
	panic("implement me")
}
