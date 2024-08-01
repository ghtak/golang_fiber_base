package main

import (
	"github.com/gofiber/fiber/v3"
	"github.com/golang_fiber_base/internal/app"
	"github.com/golang_fiber_base/internal/log"
	"github.com/golang_fiber_base/internal/module/simple"
	"go.uber.org/zap"
	"time"
)

func main() {
	logConfig := log.Config{
		Level:       "DEBUG",
		RollingFile: "./logs/app.log",
	}
	fiberConfig := fiber.Config{}
	app_ := app.NewApp(logConfig, fiberConfig)

	app_.UseFromLocals()
	app_.Logger.Info("Info Logs")
	app_.Logger.Debug("Debug Logs",
		// Structured context as strongly typed Field values.
		zap.String("string", "stringvalue"),
		zap.Int("int", 3),
		zap.Duration("duration", time.Second),
	)

	simpleModule := app_.Fiber.Group("/simple")
	simple.RegisterRoute(simpleModule)

	err := app_.Listen(":3003")
	app_.Logger.Info(err.Error())
}
