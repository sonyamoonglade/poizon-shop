package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"dto"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrNoTitle      = errors.New("missing title")
	ErrNoCategoryID = errors.New("missing categoryId")
)

func (h *Handler) GetCategoryByID(c *fiber.Ctx) error {
	categoryID, err := h.getCategoryIdFromParams(c)
	if err != nil {
		return fmt.Errorf("get category id from params: %w", err)
	}
	category, err := h.categoryService.GetByID(c.Context(), categoryID)
	if err != nil {
		return fmt.Errorf("get by id: %w", err)
	}
	return c.Status(http.StatusOK).JSON(category)
}

func (h *Handler) NewCategory(c *fiber.Ctx) error {
	title := c.Query("title", "")
	inStock := c.QueryBool("inStock", false)
	t, err := url.QueryUnescape(title)
	if err != nil {
		return fmt.Errorf("query unescape: %w", err)
	}
	title = t
	if title == "" {
		return ErrNoTitle
	}
	if err := h.categoryService.New(c.Context(), title, inStock); err != nil {
		return fmt.Errorf("new category: %w", err)
	}
	return c.SendStatus(http.StatusCreated)
}

func (h *Handler) GetAllCategories(c *fiber.Ctx) error {
	inStock := c.QueryBool("inStock", false)
	categories, err := h.categoryService.GetAllByInStock(c.Context(), inStock)
	if err != nil {
		return fmt.Errorf("get all: %w", err)
	}
	return c.Status(http.StatusOK).JSON(categories)
}

func (h *Handler) UpdateCategory(c *fiber.Ctx) error {
	var inp dto.UpdateCategoryDTO
	if err := c.BodyParser(&inp); err != nil {
		return fmt.Errorf("body parser: %w", err)
	}
	id, err := h.getCategoryIdFromParams(c)
	if err != nil {
		return fmt.Errorf("get category from params: %w", err)
	}
	if err := h.categoryService.Update(c.Context(), id, inp); err != nil {
		return fmt.Errorf("update category: %w", err)
	}
	return c.SendStatus(http.StatusOK)
}

func (h *Handler) DeleteCategory(c *fiber.Ctx) error {
	id, err := h.getCategoryIdFromParams(c)
	if err != nil {
		return fmt.Errorf("get category from params: %w", err)
	}
	if err := h.categoryService.Delete(c.Context(), id); err != nil {
		return fmt.Errorf("delete category: %w", err)
	}
	return c.SendStatus(http.StatusOK)
}

func (h *Handler) getCategoryIdFromParams(c *fiber.Ctx) (primitive.ObjectID, error) {
	categoryID := c.Params("categoryId", "")
	if categoryID == "" {
		return primitive.ObjectID{}, ErrNoCategoryID
	}

	objId, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		return primitive.ObjectID{}, fmt.Errorf("new object id: %w", err)
	}
	return objId, nil
}
