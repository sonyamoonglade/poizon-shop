package services

import (
	"context"
	"fmt"

	"domain"
	"dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"repositories"
)

type householdCategoryService struct {
	repo repositories.HouseholdCategory
}

func NewHouseholdCategoryService(repo repositories.HouseholdCategory) *householdCategoryService {
	return &householdCategoryService{
		repo: repo,
	}
}

func (h *householdCategoryService) New(ctx context.Context, title string) error {
	category := domain.NewHouseholdCategory(title)
	rank, err := h.repo.GetTopRank(ctx)
	if err != nil {
		return fmt.Errorf("get top rank: %w", err)
	}
	category.SetRank(rank)
	if err := h.repo.New(ctx, category); err != nil {
		return fmt.Errorf("new: %w", err)
	}
	return nil
}

func (h *householdCategoryService) Delete(ctx context.Context, categoryID primitive.ObjectID) error {
	return h.repo.Delete(ctx, categoryID)
}

func (h *householdCategoryService) Update(ctx context.Context, categoryID primitive.ObjectID, dto dto.UpdateCategoryDTO) error {
	return h.repo.Update(ctx, categoryID, dto)
}

func (h *householdCategoryService) GetAll(ctx context.Context) ([]domain.Category, error) {
	return h.repo.GetAll(ctx)
}
