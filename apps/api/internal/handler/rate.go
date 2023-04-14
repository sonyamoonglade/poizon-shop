package handler

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) UpdateRate(c *fiber.Ctx) error {
	newRate := c.QueryFloat("rate", 0.0)
	if newRate == 0.0 {
		return fmt.Errorf("empty rate")
	}
	if err := h.rateProvider.UpdateRate(c.Context(), newRate); err != nil {
		return fmt.Errorf("update rate: %w", err)
	}
	return c.SendStatus(http.StatusOK)
}

func (h *Handler) CurrentRate(c *fiber.Ctx) error {
	rate, err := h.rateProvider.GetYuanRate(c.Context())
	if err != nil {
		return fmt.Errorf("get yuan rate: %w", err)
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"rate": rate,
	})
}
