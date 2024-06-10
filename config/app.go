package config

import (
	"car-parking-api/config/redis"
	health "car-parking-api/internal/health/delivery/http"
	parking "car-parking-api/internal/parking/delivery/http"
	parkingRepoMemory "car-parking-api/internal/parking/repository/memory"
	parkingUseCase "car-parking-api/internal/parking/usecase"
	"car-parking-api/internal/route"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type BootstrapConfig struct {
	App                *fiber.App
	DB                 *mongo.Database
	Log                *zap.Logger
	Config             *viper.Viper
	RedisCache         *redis.Client
	PassportMiddleware fiber.Handler
}

func Bootstrap(config *BootstrapConfig) {

	healthController := health.NewHealthController()

	carRepo := parkingRepoMemory.NewCarRepository()
	parkingLotRepo := parkingRepoMemory.NewParkingSlotRepository()
	ticketRepo := parkingRepoMemory.NewTicketRepository()

	parkUseCase := parkingUseCase.NewParkingUsecase(carRepo, parkingLotRepo, ticketRepo)
	parkingController := parking.NewParkingController(parkUseCase)

	routeConfig := route.ConfigRoute{
		App:               config.App,
		HealthController:  healthController,
		ParkingController: parkingController,
	}
	routeConfig.Setup()
}
