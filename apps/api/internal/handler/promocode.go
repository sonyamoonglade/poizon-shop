package handler

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"dto"
	"go.uber.org/multierr"
	"onlineshop/api/internal/input"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrNoPromocodeID = errors.New("missing promocodeId")
)

func (h *Handler) NewPromocode(c *fiber.Ctx) error {
	var inp input.NewPromocodeInput
	if err := c.BodyParser(&inp); err != nil {
		return fmt.Errorf("body parser: %w", err)
	}

	if err := h.repositories.Promocode.Save(c.Context(), inp.ToDomainPromocode()); err != nil {
		return fmt.Errorf("save: %w", err)
	}

	return c.SendStatus(http.StatusCreated)
}

func (h *Handler) GetAllPromocodes(c *fiber.Ctx) error {
	promocodes, err := h.repositories.Promocode.GetAll(c.Context())
	if err != nil {
		return fmt.Errorf("get all: %w", err)
	}
	return c.Status(http.StatusOK).JSON(promocodes)
}

func (h *Handler) DeletePromocode(c *fiber.Ctx) error {
	promoId, err := h.getPromocodeIdFromParams(c)
	if err != nil {
		return fmt.Errorf("get promo id from params: %w", err)
	}
	if err := h.repositories.Promocode.Delete(c.Context(), promoId); err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	hCustomers, err := h.services.HouseholdCustomer.GetAllByPromocodeID(c.Context(), promoId)
	if err != nil {
		return fmt.Errorf("get all by promoid: %w", err)
	}
	cCustomers, err := h.services.ClothingCustomer.GetAllByPromocodeID(c.Context(), promoId)
	if err != nil {
		return fmt.Errorf("get all by promoid: %w", err)
	}
	// clothing
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go func() {
		for _, customer := range cCustomers {
			err = multierr.Append(err, h.services.ClothingCustomer.Update(c.Context(), customer.CustomerID, dto.UpdateClothingCustomerDTO{
				PromocodeID: &primitive.ObjectID{},
			}))
		}
		wg.Done()
	}()
	// household
	go func() {
		for _, customer := range hCustomers {
			err = multierr.Append(err, h.services.HouseholdCustomer.Update(c.Context(), customer.CustomerID, dto.UpdateHouseholdCustomerDTO{
				PromocodeID: &primitive.ObjectID{},
			}))
		}
		wg.Done()
	}()
	wg.Wait()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(http.StatusOK)
}

func (h *Handler) GetByID(c *fiber.Ctx) error {
	promoId, err := h.getPromocodeIdFromParams(c)
	if err != nil {
		return fmt.Errorf("get promo id from params: %w", err)
	}
	promo, err := h.repositories.Promocode.GetByID(c.Context(), promoId)
	if err != nil {
		return fmt.Errorf("get by id: %w", err)
	}
	return c.Status(http.StatusOK).JSON(promo)
}

func (h *Handler) getPromocodeIdFromParams(c *fiber.Ctx) (primitive.ObjectID, error) {
	promocodeId := c.Params("promocodeId", "")
	if promocodeId == "" {
		return primitive.ObjectID{}, ErrNoPromocodeID
	}

	objId, err := primitive.ObjectIDFromHex(promocodeId)
	if err != nil {
		return primitive.ObjectID{}, fmt.Errorf("new object id: %w", err)
	}
	return objId, nil
}
