package middleware

import (
	"car-parking-api/pkg/logger"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// Define excluded endpoints
var excludedEndpoints = []string{
	"/vgo/lb-status",
	// Add more endpoints as needed
}

func shouldSkipLogging(path string) bool {
	for _, endpoint := range excludedEndpoints {
		if path == endpoint {
			return true
		}
	}
	return false
}

func RamUsageMiddleware(c *fiber.Ctx) error {
	startTime := time.Now()

	// Execute next middleware(s)
	err := c.Next()

	if !shouldSkipLogging(c.Path()) {
		// Log RAM usage
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		// Log request duration
		duration := time.Since(startTime)
		logger.WithContext(c).Info(
			"ram_usage",
			zap.String("Path", c.Path()),
			zap.Any("duration", duration),
			zap.Any("MiB", m.Alloc/1024/1024),
		)
	}

	return err
}
