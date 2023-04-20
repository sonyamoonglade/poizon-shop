package handler

import (
	"context"
	"net/http"

	"onlineshop/api/internal/auth"
	"repositories"
	"services"

	"github.com/gofiber/fiber/v2"
)

type RateProvider interface {
	UpdateRate(ctx context.Context, rate float64) error
	GetYuanRate(ctx context.Context) (float64, error)
}

type Handler struct {
	repositories repositories.Repositories
	services     services.Services
	rateProvider RateProvider
}

func NewHandler(repos repositories.Repositories, services services.Services) *Handler {
	return &Handler{
		repositories: repos,
		rateProvider: repos.Rate,
		services:     services,
	}
}

func (h *Handler) RegisterRoutes(router fiber.Router, apiKey string) {
	router.Get("/", h.Home)

	api := router.Group("/api")
	api.Use(auth.NewAPIKeyMiddleware(apiKey))

	api.Post("/updateRate", h.UpdateRate)
	api.Get("/currentRate", h.CurrentRate)

	promocode := api.Group("/promocode")
	{
		promocode.Get("/all", h.GetAllPromocodes)
		promocode.Get("/:promocodeId", h.GetByID)
		promocode.Post("/new", h.NewPromocode)
		promocode.Post("/:promocodeId/delete", h.DeletePromocode)
	}

	order := api.Group("/order/:source")
	{
		order.Get("/all", h.GetAllOrders)
		order.Get("/:shortId", h.GetOrderByID)
		order.Put("/:orderId/addComment", h.AddCommentToOrder)
		order.Put("/:orderId/changeStatus", h.ChangeOrderStatus)
		order.Put("/:orderId/approve", h.Approve)
		order.Post("/:orderId/delete/", h.Delete)
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
		householdCategories.Get("/:categoryId", h.GetCategoryByID)
		householdCategories.Post("/new", h.NewCategory)
		householdCategories.Put("/:categoryId/update", h.UpdateCategory)
		householdCategories.Post("/:categoryId/delete", h.DeleteCategory)
	}

}
func (h *Handler) Home(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusOK)
}
