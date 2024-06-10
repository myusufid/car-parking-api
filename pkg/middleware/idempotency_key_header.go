package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func IdempotencyRequired(c *fiber.Ctx) error {
	// Get the token from the Authorization header
	keyString := c.Get("X-Idempotency-Key")

	// Check if the token is missing
	if keyString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Missing Idempotency Key",
		})
	}

	// Continue with the next handler
	return c.Next()
}
