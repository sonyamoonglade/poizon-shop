package handler

import (
	"errors"
	"fmt"
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

	return nil
}

func (h *Handler) GetAllPromocodes(c *fiber.Ctx) error {
	return nil
}

func (h *Handler) DeletePromocode(c *fiber.Ctx) error {
	return nil
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
