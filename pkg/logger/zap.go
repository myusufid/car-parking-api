package logger

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func WithContext(c *fiber.Ctx) *zap.Logger {
	return c.Locals("logger").(*zap.Logger)
}
