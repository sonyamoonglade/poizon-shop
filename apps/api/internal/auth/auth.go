package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

const (
	headerApiKey = "X-Api-Key"
)

var (
	ErrNoApiKey      = errors.New("missing api key")
	ErrInvalidApiKey = errors.New("invalid api key")
)

func NewAPIKeyMiddleware(key string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		headers := c.GetReqHeaders()
		candidate, ok := headers[headerApiKey]

		if !ok {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": ErrNoApiKey.Error(),
			})
		}
		if strings.TrimSpace(candidate) != key {
			return c.Status(http.StatusForbidden).JSON(fiber.Map{
				"error": ErrInvalidApiKey.Error(),
			})
		}

		return c.Next()
	}
}
