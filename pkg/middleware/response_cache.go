package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/patrickmn/go-cache"
	"time"
)

var (
	cachingClient = cache.New(10*time.Minute, 5*time.Minute)
)

func CacheMiddleware(c *fiber.Ctx) error {
	// Only cache GET requests
	if c.Method() != "GET" {
		return c.Next()
	}

	if c.Query("noCache") == "true" {
		return c.Next()
	}

	// Create a cache key based on the full URL path and IP address
	cacheKey := fmt.Sprintf("%s_%s", c.OriginalURL(), c.IP())

	// Try to retrieve the cached response
	if cachedResponse, found := cachingClient.Get(cacheKey); found {
		// Unmarshal the cached response into a JSON object
		var jsonResponse interface{}
		err := json.Unmarshal(cachedResponse.([]byte), &jsonResponse)
		if err != nil {
			// Handle error if JSON decoding fails
			return err
		}

		// Return the JSON response
		return c.JSON(jsonResponse)
	}

	// If not found, proceed with the request and cache the response
	if err := c.Next(); err != nil {
		return err
	}

	// Cache the response for 10 minutes
	cachingClient.Set(cacheKey, c.Response().Body(), cache.DefaultExpiration)

	return nil
}

// ClearCacheHandler Handler to clear the entire cache
func ClearCacheHandler(c *fiber.Ctx) error {
	cachingClient.Flush()
	return c.SendString("Cache cleared!")
}
