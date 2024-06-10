package route

import (
	health "car-parking-api/internal/health/delivery/http"
	parking "car-parking-api/internal/parking/delivery/http"
	"github.com/gofiber/fiber/v2"
)

type ConfigRoute struct {
	App               *fiber.App
	HealthController  *health.HealthController
	ParkingController *parking.ParkingController
}

func (c *ConfigRoute) Setup() {
	c.App.Get("/", c.HealthController.HelloWorld)
	c.App.Get("health", c.HealthController.GetHealth)
	c.App.Post("register", c.ParkingController.RegisterCar)
	c.App.Post("exit", c.ParkingController.ExitCar)
	c.App.Get("total_car", c.ParkingController.GetCarCountByType)
	c.App.Get("license_by_color", c.ParkingController.GetCarsByColor)
	c.v1Route()
	c.handle404()
}

func (c *ConfigRoute) handle404() {
	c.App.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "404 Not Found",
		})
	})
}

func (c *ConfigRoute) v1Route() {
	v1 := c.App.Group("v1")
	v1.Get("health", c.HealthController.GetHealth)
}
