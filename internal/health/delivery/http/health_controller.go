package http

import "github.com/gofiber/fiber/v2"

type HealthController struct {
}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (controller *HealthController) GetHealth(ctx *fiber.Ctx) error {
	return ctx.JSON("car-parking-api (live)")
}

func (controller *HealthController) HelloWorld(ctx *fiber.Ctx) error {
	return ctx.JSON("car-parking-api")
}
