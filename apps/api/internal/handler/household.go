package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"onlineshop/api/internal/input"
	"utils/boolconv"
	"utils/transliterators"
)

func (h *Handler) CallbackQueryCalculator(c *fiber.Ctx) error {
	var inp input.CallbackCalculatorQueryInput
	if err := c.QueryParser(&inp); err != nil {
		return fmt.Errorf("query parser: %w", err)
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"result": householdInjectFunc(inp.Category, inp.Subcategory, inp.ProductName),
	})
}

func householdInjectFunc(values ...string) int {
	// append random inStock value
	values = append(values, boolconv.Optimized(true))
	return len("1;" + strings.Join(transliterators.Encode(values), ";"))
}
