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

	// clothing
	wg := new(sync.WaitGroup)
	wg.Add(2)
	var deletionError error
	go func() {
		cCustomers, err := h.services.ClothingCustomer.GetAllByPromocodeID(c.Context(), promoId)
		if err != nil {
			deletionError = multierr.Append(deletionError, fmt.Errorf("get clothing customers by id: %w", err))
		}
		for _, customer := range cCustomers {
			deletionError = multierr.Append(deletionError, h.services.ClothingCustomer.Update(c.Context(), customer.CustomerID, dto.UpdateClothingCustomerDTO{
				PromocodeID: &primitive.ObjectID{},
			}))
		}
		wg.Done()
	}()
	// household
	go func() {
		hCustomers, err := h.services.HouseholdCustomer.GetAllByPromocodeID(c.Context(), promoId)
		if err != nil {
			deletionError = multierr.Append(deletionError, fmt.Errorf("get household customers by id: %w", err))
		}
		for _, customer := range hCustomers {
			deletionError = multierr.Append(deletionError, h.services.HouseholdCustomer.Update(c.Context(), customer.CustomerID, dto.UpdateHouseholdCustomerDTO{
				PromocodeID: &primitive.ObjectID{},
			}))
		}
		wg.Done()
	}()
	wg.Wait()
	if deletionError != nil {
		return fmt.Errorf("deletion error: %w", deletionError)
	}
	// after successful promocode deletion from each customer delete promocode itself
	if err := h.repositories.Promocode.Delete(c.Context(), promoId); err != nil {
		return fmt.Errorf("promocode delete: %w", err)
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

func (h *Handler) Update(c *fiber.Ctx) error {
	var inp input.UpdatePromocodeInput
	if err := c.BodyParser(&inp); err != nil {
		return fmt.Errorf("body parser: %w", err)
	}

	promocodeID, err := h.getPromocodeIdFromParams(c)
	if err != nil {
		return fmt.Errorf("get promocode id from params: %w", err)
	}

	if err := h.repositories.Promocode.Update(c.Context(), promocodeID, inp.ToDTO()); err != nil {
		return fmt.Errorf("promocode update: %w", err)
	}

	return c.SendStatus(http.StatusOK)
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
