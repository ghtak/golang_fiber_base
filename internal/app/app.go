package app

import (
	"context"
	"github.com/gofiber/fiber/v3"
	"github.com/golang_fiber_base/internal/log"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type App struct {
	Fiber  *fiber.App
	Logger *zap.Logger
}

func NewApp(logConfig log.Config, fiberConfig fiber.Config) *App {
	return &App{
		Fiber:  fiber.New(fiberConfig),
		Logger: log.NewZapLogger(logConfig),
	}
}

func (app *App) Listen(port string) error {
	return app.Fiber.Listen(port)
}

func (app *App) UseFromLocals() {
	app.Fiber.Use(func(c fiber.Ctx) error {
		app.SetLocals(c)
		return c.Next()
	})
}

func (app *App) SetLocals(c fiber.Ctx) {
	c.Locals("App", app)
}

func FromLocals(c fiber.Ctx) *App {
	return c.Locals("App").(*App)
}

func NewFiber(lc fx.Lifecycle) *fiber.App {
	app := fiber.New()
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go app.Listen(":3003")
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return app.Shutdown()
			},
		})
	return app
}
