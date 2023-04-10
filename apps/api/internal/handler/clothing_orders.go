package handler

import (
	"errors"
	"fmt"
	"net/http"

	"domain"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"onlineshop/api/internal/input"
)

var (
	ErrNoOrderID = errors.New("missing orderId")
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
	newOrder, err := h.repositories.ClothingOrder.AddComment(c.Context(), inp.ToDTO(orderID))
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
	newOrder, err := h.repositories.ClothingOrder.ChangeStatus(c.Context(), inp.ToDTO(orderID))
	if err != nil {
		return fmt.Errorf("can't change status: %w", err)
	}

	return c.Status(http.StatusOK).JSON(newOrder)
}

func (h *Handler) GetAllOrders(c *fiber.Ctx) error {
	orders, err := h.repositories.ClothingOrder.GetAll(c.Context())
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

	order, err := h.repositories.ClothingOrder.GetByShortID(c.Context(), shortId)
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
	order, err := h.repositories.ClothingOrder.Approve(c.Context(), orderID)
	if err != nil {
		return fmt.Errorf("approve: %w", err)
	}
	return c.Status(http.StatusOK).JSON(order)
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	orderID, err := h.getOrderIdFromParams(c)
	if err != nil {
		return fmt.Errorf("get order id from params: %w", err)
	}
	if err := h.repositories.ClothingOrder.Delete(c.Context(), orderID); err != nil {
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
