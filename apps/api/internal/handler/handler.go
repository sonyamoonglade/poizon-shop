package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"repositories"
	"services"
)

type RateProvider interface {
	UpdateRate(ctx context.Context, rate float64) error
	GetYuanRate(ctx context.Context) (float64, error)
}

type Handler struct {
	repositories    repositories.Repositories
	categoryService services.HouseholdCategory
	rateProvider    RateProvider
}

func NewHandler(repos repositories.Repositories, categoryService services.HouseholdCategory) *Handler {
	return &Handler{
		repositories:    repos,
		rateProvider:    repos.Rate,
		categoryService: categoryService,
	}
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	router.Get("/", h.Home)

	api := router.Group("/api")
	api.Post("/updateRate", h.UpdateRate)
	api.Get("/currentRate", h.CurrentRate)

	order := api.Group("/order/:source")
	{
		order.Get("/all", h.GetAllOrders)
		order.Get("/:shortId", h.GetOrderByID)
		order.Put("/addComment", h.AddCommentToOrder)
		order.Put("/changeStatus", h.ChangeOrderStatus)
		order.Put("/approve/:orderId", h.Approve)
		order.Post("/delete/:orderId", h.Delete)
	}

	clothingCatalog := api.Group("/clothing/catalog")
	{
		clothingCatalog.Get("/all", h.GetAllClothingCatalog)
		clothingCatalog.Post("/addItem", h.AddNewClothingProduct)
		clothingCatalog.Post("/deleteItem", h.DeleteClothingItem)
	}

	householdCategories := api.Group("/household/categories")
	{
		householdCategories.Get("/all", h.GetAllCategories)
		householdCategories.Post("/new", h.NewCategory)
		householdCategories.Put("/:categoryId/update", h.UpdateCategory)
		householdCategories.Post("/:categoryId/delete", h.DeleteCategory)
	}
}
func (h *Handler) Home(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusOK)
}

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
