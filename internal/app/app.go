package app

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/golang_fiber_base/internal/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewFiber(
	lc fx.Lifecycle,
	config config.Config,
	log *zap.Logger,
) *fiber.App {
	app := fiber.New()
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				log.Info("Server Listen", zap.String("Port", config.AppPort))
				go app.Listen(fmt.Sprintf(":%s", config.AppPort))
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return app.Shutdown()
			},
		})
	return app
}

func Module() fx.Option {
	return fx.Module(
		"ModuleApp",
		fx.Provide(NewFiber),
		fx.Invoke(func(*fiber.App) {}),
	)
}
