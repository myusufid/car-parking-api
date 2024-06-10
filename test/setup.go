package test

import (
	"car-parking-api/config"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func SetupIntegration(params ...interface{}) (*fiber.App, *viper.Viper, *mongo.Database, *zap.Logger) {
	var app *fiber.App
	var configViper *viper.Viper
	var db *mongo.Database
	var log *zap.Logger

	// Iterate over provided parameters and assign them if available
	for _, param := range params {
		switch p := param.(type) {
		case *fiber.App:
			app = p
		case *viper.Viper:
			configViper = p
		case *mongo.Database:
			db = p
		case *zap.Logger:
			log = p
		}
	}

	// If any parameter is not provided, initialize it
	if app == nil {
		app = config.NewFiber()
	}
	if configViper == nil {
		configViper = config.NewViper()
	}
	if db == nil {
		db = config.NewDatabase(configViper)
	}
	if log == nil {
		log = config.NewLogger()
	}

	return app, configViper, db, log
}
