package middleware

import (
	"car-parking-api/pkg/framework"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func ZapLogger(c *fiber.Ctx) error {
	logger := framework.CreateLogger()
	// Generate a request ID (you may use a library to generate a unique ID)
	requestID := c.GetRespHeader(fiber.HeaderXRequestID)
	// Add request ID to the logger's context
	loggerWithRequestID := logger.With(zap.String("requestID", requestID))
	c.Locals("logger", loggerWithRequestID)
	// Continue with the request
	return c.Next()
}
