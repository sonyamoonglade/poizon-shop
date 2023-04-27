package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"domain"
	"dto"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"onlineshop/api/internal/input"
)

var (
	ErrNoOrderID = errors.New("missing orderId")
	ErrNoSource  = errors.New("missing source")
)

func (h *Handler) AddCommentToOrder(c *fiber.Ctx) error {
	var inp input.AddCommentToOrderInput
	if err := c.BodyParser(&inp); err != nil {
		return fmt.Errorf("body parsing error: %w", err)
	}

	orderID, err := h.getOrderIdFromParams(c)
	if err != nil {
		return fmt.Errorf("get order id from params: %w", err)
	}
	source, err := h.getSourceFromParams(c)
	if err != nil {
		return fmt.Errorf("get source from params: %w", err)
	}
	if source == domain.SourceClothing {
		newOrder, err := h.repositories.ClothingOrder.AddComment(c.Context(), inp.ToDTO(orderID))
		if err != nil {
			return fmt.Errorf("can't add comment: %w", err)
		}

		return c.Status(http.StatusOK).JSON(newOrder)
	}
	newOrder, err := h.repositories.HouseholdOrder.AddComment(c.Context(), inp.ToDTO(orderID))
	if err != nil {
		return fmt.Errorf("can't add comment: %w", err)
	}

	return c.Status(http.StatusOK).JSON(newOrder)

}

func (h *Handler) ChangeOrderStatus(c *fiber.Ctx) error {
	var inp input.ChangeOrderStatusInput
	if err := c.BodyParser(&inp); err != nil {
		return err
	}
	if ok := domain.IsValidOrderStatus(domain.Status(inp.NewStatus)); !ok {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid status value",
		})
	}
	orderID, err := h.getOrderIdFromParams(c)
	if err != nil {
		return fmt.Errorf("get order id from params: %w", err)
	}
	source, err := h.getSourceFromParams(c)
	if err != nil {
		return fmt.Errorf("get source from params: %w", err)
	}

	if source == domain.SourceClothing {
		newOrder, err := h.repositories.ClothingOrder.ChangeStatus(c.Context(), inp.ToDTO(orderID))
		if err != nil {
			return fmt.Errorf("can't change status: %w", err)
		}

		return c.Status(http.StatusOK).JSON(newOrder)
	}

	newOrder, err := h.repositories.HouseholdOrder.ChangeStatus(c.Context(), inp.ToDTO(orderID))
	if err != nil {
		return fmt.Errorf("can't change status: %w", err)
	}

	return c.Status(http.StatusOK).JSON(newOrder)
}

func (h *Handler) GetAllOrders(c *fiber.Ctx) error {
	source, err := h.getSourceFromParams(c)
	if err != nil {
		return fmt.Errorf("get source from params: %w", err)
	}

	if source == domain.SourceClothing {
		orders, err := h.repositories.ClothingOrder.GetAll(c.Context())
		if err != nil {
			return fmt.Errorf("get all orders: %w", err)
		}

		return c.Status(http.StatusOK).JSON(orders)
	}
	orders, err := h.repositories.HouseholdOrder.GetAll(c.Context())
	if err != nil {
		return fmt.Errorf("get all orders: %w", err)
	}

	return c.Status(http.StatusOK).JSON(orders)
}

func (h *Handler) GetOrderByID(c *fiber.Ctx) error {
	shortId := c.Params("shortId", "")
	if shortId == "" {
		return fmt.Errorf("invalid shortId")
	}

	source, err := h.getSourceFromParams(c)
	if err != nil {
		return fmt.Errorf("get source from params: %w", err)
	}

	if source == domain.SourceClothing {
		order, err := h.repositories.ClothingOrder.GetByShortID(c.Context(), shortId)
		if err != nil {
			return fmt.Errorf("get by short id: %w", err)
		}

		return c.Status(http.StatusOK).JSON(order)
	}

	order, err := h.repositories.HouseholdOrder.GetByShortID(c.Context(), shortId)
	if err != nil {
		return fmt.Errorf("get by short id: %w", err)
	}

	return c.Status(http.StatusOK).JSON(order)
}

func (h *Handler) Approve(c *fiber.Ctx) error {
	orderID, err := h.getOrderIdFromParams(c)
	if err != nil {
		return fmt.Errorf("get order id from params: %w", err)
	}
	source, err := h.getSourceFromParams(c)
	if err != nil {
		return fmt.Errorf("get source from params: %w", err)
	}
	// add usecase with promo
	//todo::!!
	if source == domain.SourceHousehold {
		order, err := h.repositories.HouseholdOrder.Approve(c.Context(), orderID)
		if err != nil {
			return fmt.Errorf("household approve: %w", err)
		}

		if err := h.updateHouseholdPromocodeCounters(c.Context(), order); err != nil {
			return fmt.Errorf("update household promo counters: %w", err)
		}

		return c.Status(http.StatusOK).JSON(order)
	}

	order, err := h.repositories.ClothingOrder.Approve(c.Context(), orderID)
	if err != nil {
		return fmt.Errorf("clothing approve: %w", err)
	}
	if err := h.updateClothingPromocodeCounters(c.Context(), order); err != nil {
		return fmt.Errorf("update household promo counters: %w", err)
	}

	return c.Status(http.StatusOK).JSON(order)
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	orderID, err := h.getOrderIdFromParams(c)
	if err != nil {
		return fmt.Errorf("get order id from params: %w", err)
	}

	source, err := h.getSourceFromParams(c)
	if err != nil {
		return fmt.Errorf("get source from params: %w", err)
	}

	if source == domain.SourceClothing {
		if err := h.repositories.ClothingOrder.Delete(c.Context(), orderID); err != nil {
			return fmt.Errorf("delete: %w", err)
		}

		return c.SendStatus(http.StatusOK)
	}

	if err := h.repositories.HouseholdOrder.Delete(c.Context(), orderID); err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return c.SendStatus(http.StatusOK)

}

func (h *Handler) getOrderIdFromParams(c *fiber.Ctx) (primitive.ObjectID, error) {
	orderID := c.Params("orderId", "")
	if orderID == "" {
		return primitive.ObjectID{}, ErrNoOrderID
	}

	objId, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return primitive.ObjectID{}, fmt.Errorf("new object id: %w", err)
	}
	return objId, nil
}
func (h *Handler) getSourceFromParams(c *fiber.Ctx) (domain.Source, error) {
	sourceStr := c.Params("source", "")
	if sourceStr == "" {
		return domain.SourceNone, ErrNoSource
	}

	return domain.SourceFromString(sourceStr), nil
}

// todo: remove, move to service or usecase
// hurrying
func (h *Handler) updateHouseholdPromocodeCounters(ctx context.Context, order domain.HouseholdOrder) error {
	customer, err := h.services.HouseholdCustomer.GetByTelegramID(ctx, order.Customer.TelegramID)
	if err != nil {
		return fmt.Errorf("household get by telegram id: %w", err)
	}
	// In order to join Promocode field
	order.Customer = customer

	isFirstOrder, err := h.services.HouseholdOrder.HasOnlyOneOrder(ctx, customer.CustomerID)
	if err != nil {
		return fmt.Errorf("has only one order: %w", err)
	}

	promo := customer.MustGetPromocode()
	if isFirstOrder {
		promo.IncrementAsFirst(order.Cart.Size())
	} else {
		promo.IncrementAsSecondEtc(order.Cart.Size())
	}
	err = h.repositories.Promocode.Update(ctx, promo.PromocodeID, dto.UpdatePromocodeDTO{
		Counters: &promo.Counters,
	})
	if err != nil {
		return fmt.Errorf("household promocode update: %w", err)
	}
	return nil
}
func (h *Handler) updateClothingPromocodeCounters(ctx context.Context, order domain.ClothingOrder) error {
	if !order.Customer.HasPromocode() {
		return nil
	}
	customer, err := h.services.ClothingCustomer.GetByTelegramID(ctx, order.Customer.TelegramID)
	if err != nil {
		return fmt.Errorf("household get by telegram id: %w", err)
	}
	// In order to join Promocode field
	order.Customer = customer

	isFirstOrder, err := h.services.ClothingOrder.HasOnlyOneOrder(ctx, customer.CustomerID)
	if err != nil {
		return fmt.Errorf("has only one order: %w", err)
	}

	promo := customer.MustGetPromocode()
	if isFirstOrder {
		promo.IncrementAsFirst(order.Cart.Size())
	} else {
		promo.IncrementAsSecondEtc(order.Cart.Size())
	}
	err = h.repositories.Promocode.Update(ctx, promo.PromocodeID, dto.UpdatePromocodeDTO{
		Counters: &promo.Counters,
	})
	if err != nil {
		return fmt.Errorf("clothing promocode update: %w", err)
	}
	return nil
}
