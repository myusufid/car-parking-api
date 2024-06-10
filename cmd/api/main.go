package main

import (
	"car-parking-api/config"
	"car-parking-api/config/redis"
	"fmt"
)

func main() {
	log := config.NewLogger()
	viperConfig := config.NewViper()
	db := config.NewDatabase(viperConfig)
	redisClient := redis.NewRedisClient(viperConfig)
	app := config.NewFiber(viperConfig)

	config.Bootstrap(&config.BootstrapConfig{
		Log:        log,
		DB:         db,
		App:        app,
		Config:     viperConfig,
		RedisCache: redisClient,
	})

	appPort := viperConfig.GetInt("APP_PORT")
	err := app.Listen(fmt.Sprintf(":%d", appPort))
	if err != nil {
		fmt.Println("error starting application", err)
	}
}
