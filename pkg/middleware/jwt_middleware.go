package middleware

import (
	"car-parking-api/pkg/utils/jwt"

	"github.com/gofiber/fiber/v2"
)

// JWTMiddleware is a middleware to authenticate requests using JWT
func JWTMiddleware(c *fiber.Ctx) error {
	// Get the token from the Authorization header
	tokenString := c.Get("Authorization")

	// Check if the token is missing
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Missing Authorization token",
		})
	}

	// Verify the token
	claims, err := jwt.VerifyJWT(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token",
		})
	}

	// Set the claims in the context for further use in the handlers
	c.Locals("claims", claims)

	// Continue with the next handler
	return c.Next()

}
