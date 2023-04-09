package handler

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"onlineshop/api/internal/input"
)

func (h *Handler) AddNewClothingProduct(c *fiber.Ctx) error {
	var inp input.AddItemToCatalogInput
	if err := c.BodyParser(&inp); err != nil {
		return err
	}
	catalog, err := h.repositories.ClothingCatalog.GetCatalog(c.Context())
	if err != nil {
		return fmt.Errorf("get last rank: %w", err)
	}

	rank := len(catalog)

	if err := h.repositories.ClothingCatalog.AddItem(c.Context(), inp.ToNewClothingProduct(uint(rank))); err != nil {
		return fmt.Errorf("add item to catalog: %w", err)
	}

	if err := h.repositories.ClothingCustomer.NullifyCatalogOffsets(c.Context()); err != nil {
		return fmt.Errorf("nullify customer catalog offsets: %w", err)
	}

	return h.withNewCatalog(c)
}
func (h *Handler) DeleteClothingItem(c *fiber.Ctx) error {
	var inp input.RemoveItemFromCatalogInput
	if err := c.BodyParser(&inp); err != nil {
		return err
	}

	if err := h.repositories.ClothingCatalog.RemoveItem(c.Context(), inp.ItemID); err != nil {
		return fmt.Errorf("remove item: %w", err)
	}

	if err := h.repositories.ClothingCustomer.NullifyCatalogOffsets(c.Context()); err != nil {
		return fmt.Errorf("nullify customer catalog offsets: %w", err)
	}

	return h.withNewCatalog(c)
}

func (h *Handler) withNewCatalog(c *fiber.Ctx) error {
	newCatalog, err := h.repositories.ClothingCatalog.GetCatalog(c.Context())
	if err != nil {
		return fmt.Errorf("get catalog: %w", err)
	}
	return c.Status(http.StatusOK).JSON(newCatalog)
}

func (h *Handler) GetAllClothingCatalog(c *fiber.Ctx) error {
	newCatalog, err := h.repositories.ClothingCatalog.GetCatalog(c.Context())
	if err != nil {
		return fmt.Errorf("get catalog: %w", err)
	}
	return c.Status(http.StatusOK).JSON(newCatalog)
}
