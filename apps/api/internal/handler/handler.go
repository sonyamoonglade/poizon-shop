package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"onlineshop/api/internal/auth"
	"repositories"
	"services"
)

type RateProvider interface {
	UpdateRate(ctx context.Context, rate float64) error
	GetYuanRate(ctx context.Context) (float64, error)
}

type PutFileDTO struct {
	Filename,
	Destination,
	ContentType string
	Bytes []byte
}

type ImageUploader interface {
	Put(ctx context.Context, dto PutFileDTO) error
	UrlToResource(filename string) string
}

type Handler struct {
	repositories  repositories.Repositories
	services      services.Services
	rateProvider  RateProvider
	imageUploader ImageUploader
}

func NewHandler(repos repositories.Repositories, services services.Services, imageUploader ImageUploader) *Handler {
	return &Handler{
		repositories:  repos,
		rateProvider:  repos.Rate,
		services:      services,
		imageUploader: imageUploader,
	}
}

func (h *Handler) RegisterRoutes(router fiber.Router, apiKey string) {
	router.Get("/", h.Home)

	api := router.Group("/api")
	api.Use(auth.NewAPIKeyMiddleware(apiKey))

	api.Post("/upload", h.Upload)

	api.Post("/updateRate", h.UpdateRate)
	api.Get("/currentRate", h.CurrentRate)

	promocode := api.Group("/promocodes")
	{
		promocode.Get("/all", h.GetAllPromocodes)
		promocode.Get("/:promocodeId", h.GetByID)
		promocode.Post("/:promocodeId/update", h.Update)
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
		order.Post("/:orderId/delete", h.Delete)
	}

	clothingCatalog := api.Group("/clothing/catalog")
	{
		clothingCatalog.Get("/all", h.GetAllClothingCatalog)
		clothingCatalog.Post("/addItem", h.AddNewClothingProduct)
		clothingCatalog.Post("/deleteItem", h.DeleteClothingItem)
	}

	household := api.Group("/household")
	{
		householdCategories := household.Group("/categories")
		{
			householdCategories.Get("/all", h.GetAllCategories)
			householdCategories.Get("/:categoryId", h.GetCategoryByID)
			householdCategories.Post("/new", h.NewCategory)
			householdCategories.Put("/:categoryId/update", h.UpdateCategory)
			householdCategories.Post("/:categoryId/delete", h.DeleteCategory)
		}

		household.Get("/callbackCalculator", h.CallbackQueryCalculator)
	}

}
func (h *Handler) Home(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusOK)
}

// root dir in s3 storage
const rootDir = ""

func (h *Handler) Upload(c *fiber.Ctx) error {
	fheader, err := c.FormFile("file")
	if err != nil {
		return fmt.Errorf("multipart: %w", err)
	}
	f, err := fheader.Open()
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	bits, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("read all: %w", err)
	}
	seed := strconv.FormatInt(time.Now().Unix(), 10)
	path := fmt.Sprintf("%s_%s", seed, fheader.Filename)
	err = h.imageUploader.Put(c.Context(), PutFileDTO{
		Filename:    path,
		Destination: rootDir,
		ContentType: http.DetectContentType(bits),
		Bytes:       bits,
	})
	if err != nil {
		return fmt.Errorf("uploader put: %w", err)
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"url": h.imageUploader.UrlToResource(path),
	})
}
