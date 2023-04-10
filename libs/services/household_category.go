package services

import (
	"context"
	"fmt"

	"domain"
	"dto"
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

func (h *householdCategoryService) New(ctx context.Context, title string) error {
	category := domain.NewHouseholdCategory(title)
	return h.transactor.WithTransaction(ctx, func(tx context.Context) error {
		rank, err := h.repo.GetTopRank(ctx)
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
	return h.repo.Delete(ctx, categoryID)
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
