package repositories

import (
	"context"
	"errors"
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
	categories   *mongo.Collection
	onChangeHook HouseholdOnChangeFunc
}

func NewHouseholdCategoryRepo(categories *mongo.Collection, onChangeHook HouseholdOnChangeFunc) *householdCategoryRepo {
	return &householdCategoryRepo{
		categories:   categories,
		onChangeHook: onChangeHook,
	}
}

func (r *householdCategoryRepo) GetTopRank(ctx context.Context) (uint, error) {
	numDocuments, err := r.categories.CountDocuments(ctx, bson.D{})
	if err != nil {
		return 0, fmt.Errorf("count documents: %w", err)
	}
	return uint(numDocuments), nil
}

func (r *householdCategoryRepo) Save(ctx context.Context, c domain.HouseholdCategory) error {
	_, err := r.categories.InsertOne(ctx, c)
	if err != nil {
		return fmt.Errorf("insert one: %w", err)
	}
	// todo: trigger on change hook
	return r.runOnChangeHook(ctx)
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

	res, err := r.categories.UpdateOne(ctx, fil, bson.M{"$set": upd})
	if err != nil {
		return fmt.Errorf("update one: %w", err)
	}
	if res.ModifiedCount == 0 {
		return domain.ErrCategoryNotFound
	}

	return r.runOnChangeHook(ctx)
}

func (r *householdCategoryRepo) Delete(ctx context.Context, categoryID primitive.ObjectID) error {
	res, err := r.categories.DeleteOne(ctx, bson.M{"_id": categoryID})
	if err != nil {
		return fmt.Errorf("delete one: %w", err)
	}
	if res.DeletedCount == 0 {
		return domain.ErrCategoryNotFound
	}
	return r.runOnChangeHook(ctx)
}

func (r *householdCategoryRepo) GetAll(ctx context.Context) ([]domain.HouseholdCategory, error) {
	opts := options.Find()
	opts.SetSort(bson.M{"rank": 1})
	cur, err := r.categories.Find(ctx, bson.D{}, opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrNoCategories
		}
		return nil, fmt.Errorf("find: %w", err)
	}
	categories := make([]domain.HouseholdCategory, 0)
	if err := cur.All(ctx, &categories); err != nil {
		return nil, fmt.Errorf("cursor all: %w", err)
	}
	return categories, nil
}
func (r *householdCategoryRepo) GetByID(ctx context.Context, categoryID primitive.ObjectID) (domain.HouseholdCategory, error) {
	res := r.categories.FindOne(ctx, bson.M{"_id": categoryID})
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return domain.HouseholdCategory{}, domain.ErrCategoryNotFound
		}
		return domain.HouseholdCategory{}, fmt.Errorf("find one: %w", res.Err())
	}
	var c domain.HouseholdCategory
	if err := res.Decode(&c); err != nil {
		return domain.HouseholdCategory{}, fmt.Errorf("decode: %w", err)
	}
	return c, nil
}

func (r *householdCategoryRepo) GetByTitle(ctx context.Context, title string) (domain.HouseholdCategory, error) {
	res := r.categories.FindOne(ctx, bson.M{"title": title})
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.HouseholdCategory{}, domain.ErrCategoryNotFound
		}
		return domain.HouseholdCategory{}, err
	}
	var c domain.HouseholdCategory
	if err := res.Decode(&c); err != nil {
		return domain.HouseholdCategory{}, err
	}
	return c, nil
}

func (r *householdCategoryRepo) GetProductsByCategoryAndSubcategory(ctx context.Context,
	cTitle,
	sTitle string,
	availableInStock bool) ([]domain.HouseholdProduct, error) {

	filter := bson.M{"$and": []bson.M{
		{"title": cTitle},
		{"subcategories.title": sTitle},
		{"subcategories.products.availableInStock": availableInStock},
	}}

	fields := bson.M{
		"subcategories": 1,
	}
	opts := options.FindOne().SetProjection(fields)

	res := r.categories.FindOne(ctx, filter, opts)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrProductsNotFound
		}
		return nil, err
	}
	type productsResp struct {
		Subcategories []domain.Subcategory `bson:"subcategories"`
	}
	var resp productsResp
	if err := res.Decode(&resp); err != nil {
		return nil, err
	}
	return resp.Subcategories[0].Products, nil
}

func (r *householdCategoryRepo) runOnChangeHook(ctx context.Context) error {
	if r.onChangeHook != nil {
		catalog, err := r.GetAll(ctx)
		if err != nil {
			return fmt.Errorf("get all: %w", err)
		}
		r.onChangeHook(catalog)
	}
	return nil
}
