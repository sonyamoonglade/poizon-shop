package services

import (
	"context"
	"fmt"

	"domain"
	"dto"
	"functools"
	"onlineshop/database"
	"repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (h *householdCategoryService) New(
	ctx context.Context,
	title string,
	inStock bool) error {
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
		currentCategories, err := h.repo.GetAll(ctx)
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
func (h *householdCategoryService) GetByID(ctx context.Context, categoryID primitive.ObjectID) (domain.HouseholdCategory, error) {
	return h.repo.GetByID(ctx, categoryID)
}

// tx - transaction
func (h *householdCategoryService) fixCategoriesRanks(tx context.Context, categoryID primitive.ObjectID, categories []domain.HouseholdCategory) error {
	n := len(categories)
	var foundIdx int
	return functools.ForEach(func(category domain.HouseholdCategory, i int) error {
		if category.CategoryID == categoryID {
			// Few cases here:
			foundIdx = i
			// It's first and last element, if so - just delete it
			if err := h.repo.Delete(tx, categoryID); err != nil {
				return fmt.Errorf("for each, delete: %w", err)
			}
		}
		// Don't do anything
		if n == 1 {
			return nil
		} else {
			// Starting from foundIdx reduce ranks by 1
			if i > foundIdx {
				newRank := category.Rank - 1
				if err := h.repo.Update(tx, category.CategoryID, dto.UpdateCategoryDTO{
					Rank: &newRank,
				}); err != nil {
					return fmt.Errorf("for each, update: %w", err)
				}
			}
		}

		return nil
	}, categories)
}
