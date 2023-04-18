package services

import (
	"context"
	"fmt"

	"domain"
	"dto"
	"functools"
	fn "github.com/sonyamoonglade/go_func"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"onlineshop/database"
	"repositories"
)

type householdCategoryService struct {
	repo       repositories.HouseholdCategory
	transactor database.Transactor
}

func NewHouseholdCategoryService(repo repositories.HouseholdCategory, t database.Transactor) *householdCategoryService {
	return &householdCategoryService{
		repo:       repo,
		transactor: t,
	}
}

func (h *householdCategoryService) New(ctx context.Context, title string, inStock bool) error {
	category := domain.NewHouseholdCategory(title, inStock)
	return h.transactor.WithTransaction(ctx, func(tx context.Context) error {
		rank, err := h.repo.GetTopRank(ctx, inStock)
		if err != nil {
			return fmt.Errorf("get top rank: %w", err)
		}
		category.SetRank(rank)
		if err := h.repo.Save(ctx, category); err != nil {
			return fmt.Errorf("new: %w", err)
		}
		return nil
	})
}

func (h *householdCategoryService) Delete(ctx context.Context, categoryID primitive.ObjectID) error {
	return h.transactor.WithTransaction(ctx, func(tx context.Context) error {
		category, err := h.repo.GetByID(ctx, categoryID)
		if err != nil {
			return fmt.Errorf("get by id: %w", err)
		}
		currentCategories, err := h.repo.GetAllByInStock(ctx, category.InStock)
		if err != nil {
			return fmt.Errorf("get all: %w", err)
		}
		return h.fixCategoriesRanks(tx, categoryID, currentCategories)
	})
}

func (h *householdCategoryService) Update(ctx context.Context, categoryID primitive.ObjectID, dto dto.UpdateCategoryDTO) error {
	return h.repo.Update(ctx, categoryID, dto)
}

func (h *householdCategoryService) GetAll(ctx context.Context) ([]domain.HouseholdCategory, error) {
	return h.repo.GetAll(ctx)
}
func (h *householdCategoryService) GetAllByInStock(ctx context.Context, inStock bool) ([]domain.HouseholdCategory, error) {
	return h.repo.GetAllByInStock(ctx, inStock)
}

func (h *householdCategoryService) GetByID(ctx context.Context, categoryID primitive.ObjectID) (domain.HouseholdCategory, error) {
	return h.repo.GetByID(ctx, categoryID)
}

// tx - transaction
func (h *householdCategoryService) fixCategoriesRanks(tx context.Context, categoryID primitive.ObjectID, categories []domain.HouseholdCategory) error {
	categoriesIter := fn.Of(categories)
	idx := categoriesIter.IndexOf(func(c domain.HouseholdCategory) bool {
		return c.CategoryID == categoryID
	})
	categories = categoriesIter.Filter(func(c domain.HouseholdCategory) bool {
		return c.CategoryID != categoryID
	}).Result

	for i := idx; i < len(categories); i++ {
		categories[i].Rank--
		err := h.repo.Update(tx, categories[i].CategoryID, dto.UpdateCategoryDTO{
			Rank: &categories[i].Rank,
		})
		if err != nil {
			return fmt.Errorf("update: %w", err)
		}
	}

	if err := h.repo.Delete(tx, categoryID); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

func (h *householdCategoryService) CheckIfAllProductsExist(ctx context.Context, cart []domain.HouseholdProduct, inStock bool) (bool, domain.HouseholdProduct, error) {
	categories, err := h.repo.GetAllByInStock(ctx, inStock)
	if err != nil {
		return false, domain.HouseholdProduct{}, fmt.Errorf("get all by in stock: :%w", err)
	}
	// Combine all products into one big array
	allProductsUnflat := fn.Map(
		categories,
		func(c domain.HouseholdCategory, _ int) []domain.HouseholdProduct {
			reduceFunc := func(
				acc []domain.HouseholdProduct,
				el domain.Subcategory,
			) []domain.HouseholdProduct {
				return append(acc, el.Products...)
			}
			return functools.Reduce(reduceFunc, c.Subcategories, nil)
		})
	ok, missingProducts := doAllExist(cart, fn.Flat(allProductsUnflat))
	return ok, missingProducts, nil
}

func (h *householdCategoryService) GetProductsByCategoryAndSubcategory(ctx context.Context, cTitle, sTitle string, inStock bool) ([]domain.HouseholdProduct, error) {
	return h.repo.GetProductsByCategoryAndSubcategory(ctx, cTitle, sTitle, inStock)
}

func (h *householdCategoryService) GetByTitle(ctx context.Context, title string, inStock bool) (domain.HouseholdCategory, error) {
	return h.repo.GetByTitle(ctx, title, inStock)
}

func doAllExist(cart []domain.HouseholdProduct, allProducts []domain.HouseholdProduct) (bool, domain.HouseholdProduct) {
	var missing domain.HouseholdProduct
	ok := fn.
		Of(cart).
		All(func(cartProduct domain.HouseholdProduct) bool {
			foundIdx := fn.
				Of(allProducts).
				IndexOf(func(product domain.HouseholdProduct) bool {
					return cartProduct.ProductID == product.ProductID
				})

			if foundIdx == -1 {
				missing = cartProduct
				return false
			}

			return true
		})
	return ok, missing
}
