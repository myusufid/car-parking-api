package config

import (
	"car-parking-api/pkg/exception"
	"car-parking-api/pkg/framework"
	"car-parking-api/pkg/framework/request_logger"
	"car-parking-api/pkg/logger"
	"car-parking-api/pkg/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/storage/redis/v3"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"runtime"
	"runtime/debug"
)

func NewFiber(params ...interface{}) *fiber.App {
	app := fiber.New(NewFiberConfig())
	app.Use(recover.New(recover.Config{
		EnableStackTrace:  true,
		StackTraceHandler: defaultStackTraceHandler,
	}))
	app.Use(requestid.New())

	loggerZap := framework.CreateLogger()

	defer loggerZap.Sync()

	// Use Zap logger as a middleware
	app.Use(middleware.ZapLogger)

	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Authorization,Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin,X-Idempotency-Key",
		AllowOrigins:     "http://localhost",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	var viperConfig *viper.Viper

	if len(params) > 0 {
		for _, param := range params {
			switch p := param.(type) {
			case *viper.Viper:
				viperConfig = p
			}
		}
	}

	if viperConfig != nil {
		cfg := viperConfig
		store := redis.New(redis.Config{
			Host:      cfg.GetString("REDIS_HOST"),
			Port:      cfg.GetInt("REDIS_PORT"),
			Username:  "",
			Password:  "",
			Database:  0,
			Reset:     false,
			TLSConfig: nil,
			PoolSize:  10 * runtime.GOMAXPROCS(0),
		})

		app.Use(idempotency.New(idempotency.Config{Storage: store}))
	} else {
		app.Use(idempotency.New())
	}

	app.Use(request_logger.New(request_logger.Config{
		Logger:   loggerZap,
		SkipURIs: []string{"/vgo/lb-status"},
		Fields:   []string{"method", "url", "status", "queryParams", "body", "ip", "host", "path", "ua", "bytesSent", "requestId", "error", "latency", "resBody", "reqHeaders"},
		Levels:   []zapcore.Level{zapcore.ErrorLevel, zapcore.WarnLevel, zapcore.InfoLevel, zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel},
	}))

	return app
}

func defaultStackTraceHandler(ctx *fiber.Ctx, e interface{}) {
	logger.WithContext(ctx).Error("error:", zap.Any("message", e), zap.ByteString("stack", debug.Stack()))
}

func NewFiberConfig() fiber.Config {
	return fiber.Config{
		ErrorHandler: exception.ErrorHandler,
		Prefork:      false,
	}
}
